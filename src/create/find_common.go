package create

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func FindCommonDirectory(directory, resourceType string) (source string, err error) {
	for directory != "" {
		// Check if the module directory exists in the current directory
		var module []string
		switch strings.ToLower(resourceType) {
		case "s3":
			module = []string{"terraform-aws-s3-bucket"}

		case "secrets-manager":
			module = []string{"terraform-aws-secrets-manager"}

		case "iam/policies":
			module = []string{"terraform-aws-iam", "modules", "iam-policy"}

		default:
			panic("Resource type is not supported")
		}

		moduleDir := filepath.Join(directory, "common", "terragrunt-base", "modules", filepath.Join(module...))
		if fileInfo, err := os.Stat(moduleDir); err == nil && fileInfo.IsDir() {
			if !strings.Contains(fileInfo.Name(), module[len(module)-1]) {
				// This is not the correct directory, continue searching
				continue
			}

			// Construct the source path using the module directory
			relPath, _ := filepath.Rel(directory, moduleDir)
			numUpDirs := strings.Count(relPath, string(filepath.Separator))
			if len(module) > 1 && module[len(module)-2] == "modules" {
				numUpDirs--
			}
			source = strings.Repeat("../", numUpDirs+1)
			source = filepath.Join(source, "common", "terragrunt-base", "modules", filepath.Join(module...))
			source = strings.Replace(source, "/terragrunt-base", "//terragrunt-base", 1)
			return source, nil
		}

		directory = filepath.Dir(directory)
	}

	return "", errors.New("module directory not found")
}
