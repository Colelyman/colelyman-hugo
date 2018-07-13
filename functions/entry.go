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
		fmt.Println("Connected to Github")
		repo := getRepo(client)
		fmt.Println("Retrieved the repo")
		path, file, _ := writePost(entry)
		fmt.Println("Wrote the post")
		tree, err := getTree(path, file, client, repo)
		fmt.Printf("Got the tree, with tree: %v err: %s\n", tree, err)
		if err != nil {
			return "", err
		}
		err = pushCommit(client, repo, tree)
		fmt.Printf("Pushed the commit with err: %s\n", err)

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
	repoURL := strings.Split(os.ExpandEnv("$REPOSITORY_URL"), "/")
	fmt.Printf("repoURL %v\n", repoURL)
	repo, _, err := client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+branch)
	if err != nil {
		panic(err)
	}

	// var baseRef *github.Reference
	// baseRef, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+branch)
	// if err != nil {
	// 	panic(err)
	// }
	// newRef := &github.Reference{Ref: github.String("refs/heads/" + branch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	// repo, _, err = client.Git.CreateRef(ctx, sourceOwner, sourceRepo, newRef)
	// if err != nil {
	// 	panic(err)
	// }
	return repo
}

// this function adds the new file to the repo
func getTree(path string, file string, client *github.Client, repo *github.Reference) (*github.Tree, error) {
	fmt.Printf("path: %s file: %s sourceOwner: %s sourceRepo: %s ctx: %s\n", path, file, sourceOwner, sourceRepo, ctx)
	if repo == nil {
		fmt.Println("repo is nil")
	}
	fmt.Printf("SHA: %+v\n", *repo.Object)
	if client == nil {
		fmt.Println("client is nil")
	}
	if repo == nil {
		fmt.Println("repo is nil")
	}
	tree, _, err := client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *repo.Object.SHA, []github.TreeEntry{github.TreeEntry{Path: github.String(path), Type: github.String("blob"), Content: github.String(file), Mode: github.String(("100644"))}})
	fmt.Printf("getTree err: %s\n", err)

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
	buff.WriteString("categories = [\"Micro\"]\n")
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
	buff.WriteString(entry.Content + "\n")

	path := strings.Replace(entry.slug, " ", "-", -1) + ".md"

	return "site/content/micro/" + path, buff.String(), nil
}
