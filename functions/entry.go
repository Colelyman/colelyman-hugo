package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Entry struct {
	Content string
	tags    []string
	slug    string
}

var ctx context.Context

func CreateEntry(bodyValues url.Values) (string, error) {
	if _, ok := bodyValues["content"]; ok {
		fmt.Println(bodyValues)
		entry := new(Entry)
		entry.Content = bodyValues["content"][0]
		entry.tags = bodyValues["category"]
		entry.slug = bodyValues["mp-slug"][0]

		repo := getRepo(connectGitHub())
		fmt.Println(repo)
		return "/", nil
	}
	return "",
		errors.New("Content in response body is missing")
}

func connectGitHub() *github.Client {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.ExpandEnv("$GIT_API_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func getRepo(client *github.Client) *github.Reference {
	fmt.Printf("$REPOSITORY_URL: %s", os.ExpandEnv("$REPOSITORY_URL"))
	repoURL := strings.Split(os.ExpandEnv("$REPOSITORY_URL"), "/")
	fmt.Printf("repoURL %v\n", repoURL)
	// owner, repoName := repoURL[len(repoURL)-2], repoURL[len(repoURL)-1]
	// fmt.Printf("Owner: %s and repoName: %s and branch: %s\n", owner, repoName, os.ExpandEnv("$BRANCH"))
	repo, _, err := client.Git.GetRef(ctx, "Colelyman", "colelyman-hugo", "origin/"+os.ExpandEnv("$BRANCH"))
	if err != nil {
		panic(err)
	}

	return repo
}
