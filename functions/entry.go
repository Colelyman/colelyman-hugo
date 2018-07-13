package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	hashids "github.com/speps/go-hashids"
	"golang.org/x/oauth2"
)

type Entry struct {
	Content string
	tags    []string
	slug    string
}

var ctx context.Context

var sourceOwner = "Colelyman"
var authorName = "Cole Lyman"
var authorEmail = "cole@colelyman.com"
var sourceRepo = "colelyman-hugo"
var branch = "master"

func CreateEntry(bodyValues url.Values) (string, error) {
	if _, ok := bodyValues["content"]; ok {
		fmt.Println(bodyValues)
		entry := new(Entry)
		entry.Content = bodyValues["content"][0]
		if tags, ok := bodyValues["category"]; ok {
			entry.tags = tags
		} else {
			entry.tags = nil
		}
		if slug, ok := bodyValues["mp-slug"]; ok {
			entry.slug = slug[0]
		} else {
			entry.slug = generateSlug()
		}

		client := connectGitHub()
		repo := getRepo(client)
		path, file, _ := writePost(entry)
		tree, err := getTree(path, file, client, repo)
		if err != nil {
			return "", err
		}
		err = pushCommit(client, repo, tree)
		return "/", err
	}
	return "",
		errors.New("Content in response body is missing")
}

func generateSlug() string {
	hd := hashids.NewData()
	hd.Salt = "do you want to know a secret?"
	h, _ := hashids.NewWithData(hd)
	id, _ := h.EncodeHex(time.Now().String())

	return id
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
	fmt.Printf("$REPOSITORY_URL: %s, $BRANCH: %s\n", os.Getenv("REPOSITORY_URL"), os.Getenv("BRANCH"))
	repoURL := strings.Split(os.ExpandEnv("$REPOSITORY_URL"), "/")
	fmt.Printf("repoURL %v\n", repoURL)
	// owner, repoName := repoURL[len(repoURL)-2], repoURL[len(repoURL)-1]
	repo, _, err := client.Git.GetRef(ctx, sourceOwner, sourceRepo, "heads/"+branch)
	if err != nil {
		panic(err)
	}

	return repo
}

// this function adds the new file to the repo
func getTree(path string, file string, client *github.Client, repo *github.Reference) (*github.Tree, error) {
	entries := make([]github.TreeEntry, 1)
	entries = append(entries, github.TreeEntry{Path: github.String(path), Type: github.String("blob"), Content: github.String(file), Mode: github.String(("100644"))})
	tree, _, err := client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *repo.Object.SHA, entries)

	return tree, err
}

func pushCommit(client *github.Client, repo *github.Reference, tree *github.Tree) error {
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *repo.Object.SHA)
	if err != nil {
		return err
	}

	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
	message := "Added new micropub entry."
	commit := &github.Commit{Author: author, Message: &message, Tree: tree, Parents: []github.Commit{*parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
	if err != nil {
		return err
	}

	repo.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, repo, false)
	return err
}

func writePost(entry *Entry) (string, string, error) {
	var buff bytes.Buffer

	// write the front matter in toml format
	buff.WriteString("+++\n")
	buff.WriteString("title = \"" + entry.slug + "\"\n") // TODO come up with a title
	buff.WriteString("date = \"" + time.Now().String() + "\"\n")
	buff.WriteString("categories = [\"Micro\"]")
	buff.WriteString("tags = [")
	for i, tag := range entry.tags {
		buff.WriteString("\"" + tag + "\"")
		if i < len(entry.tags)-1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString("]\n")
	buff.WriteString("slug = " + entry.slug)
	buff.WriteString("+++\n")

	// write the content
	buff.WriteString(entry.Content)

	path := strings.Replace(entry.slug, " ", "-", -1) + ".md"

	return "site/content/micro/" + path, buff.String(), nil
}
