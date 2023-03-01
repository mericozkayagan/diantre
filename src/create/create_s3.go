package create

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	s3TemplateFile    []byte
	iamTemplateFile   []byte
	placeholderValues map[string]string
)

func addParameters() {
	fmt.Println(`Leave TerragruntPath null, if your directory is like terragrunt/app
	***************************************
	`)

	placeholderValues = map[string]string{
		"Environment":    "",
		"Name":           "",
		"AllowedOrigins": "",
		"PublicAccess":   "",
		"TerragruntPath": "",
	}

	for placeholder, _ := range placeholderValues {
		fmt.Printf("Enter value for %s: ", placeholder)
		var input string
		fmt.Scanln(&input)
		placeholderValues[placeholder] = input
	}

	if !(placeholderValues["Environment"] == "pg" || placeholderValues["Environment"] == "qa" || placeholderValues["Environment"] == "prod") {
		panic(placeholderValues["Environment"])
	}

	if placeholderValues["TerragruntPath"] == "" {
		placeholderValues["TerragruntPath"] = "app"
	}

}

func CreateS3() {
	addParameters()

	switch placeholderValues["PublicAccess"] {
	case "true":
		s3TemplateFile, _ = ioutil.ReadFile("templates/s3/public-template.hcl")

	case "false":
		s3TemplateFile, _ = ioutil.ReadFile("templates/s3/private-template.hcl")
	default:
		panic("Please enter true & false to PublicAccess")
	}

	s3Content := string(s3TemplateFile)
	for placeholder, value := range placeholderValues {
		s3Content = strings.Replace(s3Content, "{{ "+placeholder+" }}", value, -1)
	}

	CreateFile(placeholderValues["Name"], placeholderValues["Environment"], s3Content, "s3", placeholderValues["TerragruntPath"])

	//Create the necessary IAM policy
	iamTemplateFile, _ = ioutil.ReadFile("templates/iam/s3-template.hcl")
	iamContent := string(iamTemplateFile)

	for placeholder, value := range placeholderValues {
		iamContent = strings.Replace(iamContent, "{{ "+placeholder+" }}", value, -1)
	}
	CreateFile(placeholderValues["Name"], placeholderValues["Environment"], iamContent, "iam/policies", placeholderValues["TerragruntPath"])

	// add dependency to custom_role_arns part
	AddToCustomRoleArns(placeholderValues["Name"], placeholderValues["Environment"], "iam/role", placeholderValues["TerragruntPath"])

	// add the dependencies part
	AddDependenciesToEnd(placeholderValues["Name"], placeholderValues["Environment"], s3Content, "iam/role", placeholderValues["TerragruntPath"])

	fmt.Println("Successfully created the S3 bucket named: ", placeholderValues["Name"], "in", placeholderValues["Environment"], "environment")
}
