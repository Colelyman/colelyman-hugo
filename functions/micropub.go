package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type IndieAuthRes struct {
	Me       string `json:"me"`
	ClientId string `json:"client_id"`
	Scope    string `json:"scope"`
	Issue    int    `json:"issued_at"`
	Nonce    int    `json:"nonce"`
}

func createEntry(bodyValues url.Values) (location string, err error) {
	if content, ok := bodyValues["content"]; ok {
		fmt.Printf("Content has length: %d\n", len(content))
		fmt.Println(content)
	}

	location = "here"
	return
}

func checkAccess(token string) (bool, error) {
	// form the request to check the token
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://tokens.indieauth.com/token", nil)
	if err != nil {
		return false,
			errors.New("Error making the request for checking token access")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)
	// send the request
	res, err := client.Do(req)
	if err != nil {
		return false,
			errors.New("Error sending the request for checking token access")
	}
	defer res.Body.Close()
	// parse the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false,
			errors.New("Error parsing the response for checking token access")
	}
	var indieAuthRes = new(IndieAuthRes)
	fmt.Println(res.Body)
	fmt.Println(body)
	err = json.Unmarshal(body, &indieAuthRes)
	if err != nil {
		return false,
			errors.New("Error parsing the response into json for checking token access " + err.Error())
	}

	// verify results of the response
	if indieAuthRes.Me != "http://colelyman.com" {
		return false,
			errors.New("Me does not match")
	}
	scopes := strings.Fields(indieAuthRes.Scope)
	postPresent := false
	for _, scope := range scopes {
		if scope == "post" || scope == "create" || scope == "update" {
			postPresent = true
			break
		}
	}
	if !postPresent {
		return false,
			errors.New("Post is not present in the scope")
	}
	return true, nil
}

func checkAuthorization(bodyValues url.Values, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var token string
	// check the headers for authorization first
	if token, ok := req.Headers["authorization"]; ok {
		fmt.Println("Authorization header exists: " + token)
	} else if token, ok := bodyValues["access_token"]; ok {
		token := "Bearer " + token[0]
		fmt.Println("Access_token in body exists: " + token)
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Unauthorized, access token was not provided",
		}, errors.New("Access token was not provided")
	}

	// var err string
	if ok, err := checkAccess(token); ok {
		location, err := createEntry(bodyValues)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: 403,
				Body:       "Error occured while checking access",
			}, err
		}
		// Everything worked out!! Send the location and an OK status
		return &events.APIGatewayProxyResponse{
			StatusCode: 201,
			Headers:    map[string]string{"Location": location},
		}, nil
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Forbidden, the provided access token does not grant permission",
		}, errors.New("The provided access token does not grant permission " + err.Error())
	}
}

func checkContentType(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if contentType, ok := req.Headers["content-type"]; ok {
		if contentType == "application/x-www-form-urlencoded" {
			bodyValues, err := url.ParseQuery(req.Body)
			if err != nil {
				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Bad Request, error parsing the body of the request",
				}, errors.New("Error parsing the body of the request")
			} else if val, ok := bodyValues["h"]; !ok || val[0] != "entry" {
				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Bad Request, either there is no h value in the body or its value is not entry",
				}, errors.New("Error with the h value in the body of the request")
			}
			// proceed to check authorization
			return checkAuthorization(bodyValues, req)
		} else {
			return &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Bad Request, content-type " + contentType + " is not supported, use application/x-www-form-urlencoded",
			}, errors.New("Content-type " + contentType + " is not supported, use application/x-www-form-urlencoded")
		}
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request, content-type is not provided in the request",
		}, errors.New("Content-type is not provided in the request")
	}
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// a handler for GET requests, used for troubleshooting
	if req.HTTPMethod == "GET" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Everything is working, this is the GET request body: " + req.Body,
		}, nil
	}
	// check if the request is a post
	if req.HTTPMethod != "POST" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "The HTTP method is not allowed, make a POST request",
		}, errors.New("HTTP method is not valid")
	}
	// fmt.Println(req.Headers)
	// fmt.Println(req.Body)

	// check the content-type and proceed down the rabbit hole
	return checkContentType(req)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
