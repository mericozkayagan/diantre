// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The commitpr command utilizes go-github as a CLI tool for
// pushing files to a branch and creating a pull request from it.
// It takes an auth token as an environment variable and creates
// the commit and the PR under the account affiliated with that token.
//
// The purpose of this example is to show how to use refs, trees and commits to
// create commits and pull requests.
//
// Note, if you want to push a single file, you probably prefer to use the
// content API. An example is available here:
// https://godoc.org/github.com/google/go-github/github#example-RepositoriesService-CreateFile
//
// Note, for this to work at least 1 commit is needed, so you if you use this
// after creating a repository you might want to make sure you set `AutoInit` to
// `true`.
package commitpr

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/mericozkayagan/diantre/src/authentication"
)

var (
	sourceOwner   string
	sourceRepo    string
	commitMessage string
	commitBranch  string
	baseBranch    string
	prRepoOwner   string
	prRepo        string
	prBranch      string
	prSubject     string
	prDescription string
	sourceFiles   string

	authorName  string
	authorEmail string
)

var client *github.Client
var ctx = context.Background()

func addParameters() {
	fmt.Println("Please enter your repository name: ")
	fmt.Scanln(&sourceOwner)

	fmt.Println("Please enter your repository description (add no spaces for now please): ")
	fmt.Scanln(&sourceRepo)

	fmt.Println("Please enter your commit message: ")
	fmt.Scanln(&commitMessage)

	fmt.Println("Please enter your commit branch: ")
	fmt.Scanln(&commitBranch)

	fmt.Println("Please enter your base branch: ")
	fmt.Scanln(&baseBranch)

	fmt.Println("Please enter your merge repo owner: ")
	fmt.Scanln(&prRepoOwner)

	fmt.Println("Please enter your merge repo: ")
	fmt.Scanln(&prRepo)

	fmt.Println("Please enter your merge branch: ")
	fmt.Scanln(&prBranch)

	fmt.Println("Please enter your pull request title: ")
	fmt.Scanln(&prSubject)

	fmt.Println("Please enter your pull request description: ")
	fmt.Scanln(&prDescription)

	fmt.Println("Please enter your files with semicolons: ")
	fmt.Scanln(&sourceFiles)

	fmt.Println("Please enter your author name: ")
	fmt.Scanln(&authorName)

	fmt.Println("Please enter your author email: ")
	fmt.Scanln(&authorEmail)
}

func CommitPr() {
	flag.Parse()
	addParameters()

	if sourceOwner == "" || sourceRepo == "" || commitBranch == "" || sourceFiles == "" || authorName == "" || authorEmail == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}

	client, ctx = authentication.Auth()

	//checkout the code and replace the values with the ones you want

	ref, err := getRef()
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	if err := createPR(); err != nil {
		log.Fatalf("Error while creating the pull request: %s", err)
	}
}

func checkOutCode() {
	fmt.Println("Checking out code...")
}

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef() (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	// We consider that an error means the branch has not been found and needs to
	// be created.
	if &commitBranch == &baseBranch {
		return nil, errors.New("the commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if baseBranch == "" {
		return nil, errors.New("the `-base-branch` should not be set to an empty string when the branch specified by `-commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+baseBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, sourceOwner, sourceRepo, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ref *github.Reference) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(sourceFiles, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	return targetName, b, err
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &github.Timestamp{date}, Name: &authorName, Email: &authorEmail}
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
	return err
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func createPR() (err error) {
	if prSubject == "" {
		return errors.New("missing `-pr-title` flag; skipping PR creation")
	}

	if prRepoOwner != "" && prRepoOwner != sourceOwner {
		commitBranch = fmt.Sprintf("%s:%s", sourceOwner, commitBranch)
	} else {
		prRepoOwner = sourceOwner
	}

	if prRepo == "" {
		prRepo = sourceRepo
	}

	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &commitBranch,
		Base:                &prBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, prRepoOwner, prRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
