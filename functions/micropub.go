package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func checkAuthorization(token string) bool {
	return false
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// a handler for GET requests, used for troubleshooting
	if request.HTTPMethod == "GET" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Everything is working, this is the GET request body: " + request.Body,
		}, nil
	}
	// check if the request is a post
	if request.HTTPMethod != "POST" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "The HTTP method is not allowed, make a POST request",
		}, errors.New("HTTP method is not valid")
	}
	fmt.Println(request.Headers)
	fmt.Println(request.Body)

	// check the content-type
	if contentType, ok := request.Headers["content-type"]; ok {
		if contentType == "application/x-www-form-urlencoded" {
			bodyValues, err := url.ParseQuery(request.Body)
			if err != nil {
				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Bad Request, error parsing the body of the request",
				}, errors.New("Error parsing the body of the request")
			}
			fmt.Println(bodyValues)
		}
	}

	// TODO Validate the token via tokens.indieauth.com
	if token, ok := request.Headers["authorization"]; ok {
		fmt.Println("Authorization header exists: " + token)
		if checkAuthorization(token) {
			// process request
		} else {
			return &events.APIGatewayProxyResponse{
				StatusCode: 403,
				Body:       "Forbidden, the provided access token does not grant permission",
			}, errors.New("The provided access token does not grant permission")
		}
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Unauthorized, access token was not provided",
		}, errors.New("Access token was not provided")
	}
	// TODO Make sure that the access token is still valid
	// TODO Validate the request parameters
	loc := "http://colelyman.com"
	// the post was successfully created!
	return &events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": loc},
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
