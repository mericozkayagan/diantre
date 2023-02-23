// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The newrepo command utilizes go-github as a cli tool for
// creating new repositories. It takes an auth token as
// an environment variable and creates the new repo under
// the account affiliated with that token.
package newrepo

import (
	"flag"
	"fmt"
	"log"

	"github.com/google/go-github/v50/github"
	"github.com/mericozkayagan/diantre/src/authentication"
)

// var (
//
//	name               = flag.String("name", "", "Name of repo to create in authenticated user's GitHub account.")
//	description        = flag.String("description", "", "Description of created repo.")
//	private            = flag.Bool("private", false, "Will created repo be private.")
//	autoInit           = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
//	templateRepository = flag.String("template-repository", "", "Template repository to use when creating the repository.")
//
// )
var (
	name                   string
	description            string
	private                bool
	autoInit               bool
	owner                  string
	templateRepositoryName string
)

func init() {
	fmt.Println("Please enter your repository name: ")
	fmt.Scanln(&name)

	fmt.Println("Please enter your repository description (add no spaces for now please): ")
	fmt.Scanln(&description)

	fmt.Println("Do you want your repository to be private (true/false): ")
	fmt.Scanf("%t\n", &private)

	fmt.Println("Do you want an empty commit with a README.md (true/false): ")
	fmt.Scanf("%t\n", &autoInit)

	fmt.Println("Please enter the owner of the template repository name: ")
	fmt.Scanln(&owner)

	fmt.Println("Please give a template repository name (optional):  ")
	fmt.Scanln(&templateRepositoryName)

	fmt.Println("Please enter your organization (optional): ")
	fmt.Scanln(&templateRepositoryName)

}

func CreateNewRepo() {
	flag.Parse()

	client, ctx := authentication.Auth()

	if templateRepositoryName == "" {
		r := &github.Repository{
			Name:        &name,
			Private:     &private,
			Description: &description,
			AutoInit:    &autoInit,
			IsTemplate:  github.Bool(false),
		}

		repo, _, err := client.Repositories.Create(ctx, "", r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully created new repo: %v", repo.GetName())
	}

	if templateRepositoryName != "" {
		r := github.TemplateRepoRequest{
			Name:               &name,
			Owner:              &owner,
			IncludeAllBranches: github.Bool(true),
			Private:            &private,
		}

		repo, _, err := client.Repositories.CreateFromTemplate(ctx, owner, templateRepositoryName, &r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Successfully created new repo from template: %v", repo.GetName())
	}
}
