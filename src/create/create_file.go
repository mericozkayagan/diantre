package create

import (
	"os"
	"strings"
)

func CreateFile(name, environment, content, resourceType, terragruntPath string) {
	dirPath := "terragrunt/" + terragruntPath + "/" + environment + "/" + resourceType + "/" + name

	commonSource, err := FindCommonDirectory(dirPath, resourceType)
	content = strings.Replace(content, "{{ "+"Source"+" }}", commonSource, -1)

	err = os.MkdirAll(dirPath, 0777) // Create directory and all its parents
	if err != nil {
		panic(err)
	}

	// Create the file
	filePath := dirPath + "/terragrunt.hcl"
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		panic(err)
	}

	err = file.Sync()
	if err != nil {
		panic(err)
	}
}
