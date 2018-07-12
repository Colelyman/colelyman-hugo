package main

import (
	"errors"
	"fmt"
	"net/url"
)

type Entry struct {
	Content []string
	tags    []string
}

func CreateEntry(bodyValues url.Values) (string, error) {
	if _, ok := bodyValues["content"]; ok {
		fmt.Println(bodyValues)
		entry := new(Entry)
		entry.Content = bodyValues["content"]
		return "/", nil
	}
	return "",
		errors.New("Content in response body is missing")
}
