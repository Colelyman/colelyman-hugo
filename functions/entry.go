package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	hashids "github.com/speps/go-hashids"
)

type Entry struct {
	Content string
	tags    []string
	slug    string
	hash    string
}

func CreateEntry(bodyValues url.Values) (string, error) {
	if _, ok := bodyValues["content"]; ok {
		fmt.Println(bodyValues)
		entry := new(Entry)
		entry.Content = bodyValues["content"][0]
		entry.hash = generateHash()
		if tags, ok := bodyValues["category"]; ok {
			entry.tags = tags
		} else {
			entry.tags = nil
		}
		if slug, ok := bodyValues["mp-slug"]; ok && len(slug) > 0 && slug[0] != "" {
			entry.slug = slug[0]
		} else {
			entry.slug = entry.hash
		}
		fmt.Printf("Slug value is %s\n", entry.slug)

		// construct the post
		path, file, _ := writePost(entry)
		err := CommitEntry(path, file)
		if err != nil {
			return "", err
		}

		return "/micro/" + entry.slug, err
	}
	return "",
		errors.New("Content in response body is missing")
}

func generateHash() string {
	hd := hashids.NewData()
	hd.Salt = "do you want to know a secret?"
	h, _ := hashids.NewWithData(hd)
	id, _ := h.EncodeHex(time.Now().String())

	return id
}

func writePost(entry *Entry) (string, string, error) {
	var buff bytes.Buffer

	location, _ := time.LoadLocation("MST")
	t := time.Now().In(location).Format(time.RFC822)
	// write the front matter in toml format
	buff.WriteString("+++\n")
	buff.WriteString("title = \"#\"\n")
	buff.WriteString("date = \"" + t + "\"\n")
	buff.WriteString("categories = [\"Micro\"]\n")
	buff.WriteString("tags = [")
	for i, tag := range entry.tags {
		buff.WriteString("\"" + tag + "\"")
		if i < len(entry.tags)-1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString("]\n")
	buff.WriteString("slug = \"" + entry.slug + "\"\n")
	buff.WriteString("+++\n")

	// write the content
	buff.WriteString(entry.Content + "\n")

	path := strings.Replace(entry.slug, " ", "-", -1) + ".md"

	return "site/content/micro/" + path, buff.String(), nil
}
