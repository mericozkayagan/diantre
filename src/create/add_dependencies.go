package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func AddDependenciesToEnd(name, environment, content, directory, terragruntPath string) {

	dirPath := "terragrunt/" + terragruntPath + "/" + environment + "/" + directory
	err := os.MkdirAll(dirPath, 0777) // Create directory and all its parents
	if err != nil {
		panic(err)
	}

	// Create the file
	filePath := dirPath + "/terragrunt.hcl"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	dependencyConfig := fmt.Sprintf(`
dependency "%s" {
	config_path = "../policies/%s"
}
	`, name, name)

	_, err = file.WriteString(dependencyConfig)
	if err != nil {
		panic(err)
	}

	err = file.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully added the dependency configuration!!!")
}

func AddToCustomRoleArns(name, environment, directory, terragruntPath string) {
	dirPath := "terragrunt/" + terragruntPath + "/" + environment + "/" + directory
	err := os.MkdirAll(dirPath, 0777) // Create directory and all its parents
	if err != nil {
		panic(err)
	}

	filePath := dirPath + "/terragrunt.hcl"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Find the index of the custom_role_policy_arns list in the file content
	customRolePolicyArnsIndex := strings.Index(string(content), "custom_role_policy_arns = [")
	if customRolePolicyArnsIndex == -1 {
		panic("custom_role_policy_arns list not found in file")
	}

	// Find the end of the custom_role_policy_arns list
	customRolePolicyArnsEndIndex := strings.Index(string(content[customRolePolicyArnsIndex:]), "]")
	if customRolePolicyArnsEndIndex == -1 {
		panic("end of custom_role_policy_arns list not found in file")
	}
	customRolePolicyArnsEndIndex += customRolePolicyArnsIndex

	// Construct the new list item with the dependency
	dependencyListItem := fmt.Sprintf(`
    dependency.%s.outputs.arn,
`, name)

	// Insert the new list item before the end of the custom_role_policy_arns list
	newContent := string(content[:customRolePolicyArnsEndIndex]) + dependencyListItem + string(content[customRolePolicyArnsEndIndex:])

	// Write the new content back to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully added the dependency to the custom_role_policy_arns list!!!")
}
