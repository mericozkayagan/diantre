package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindCommonDirectory(startingDir, resourceType string) (commonDirectory string, err error) {
	// Check if the starting directory contains a "common" directory
	commonDir := filepath.Join(startingDir, "common")
	if _, err := os.Stat(commonDir); err == nil {
		// Construct the source path using the common directory
		switch strings.ToLower(resourceType) {
		case "s3":
			source := filepath.Join(commonDir, "terragrunt-base", "modules", "terraform-aws-s3-bucket")
			return source, nil
		case "secrets-manager":
			source := filepath.Join(commonDir, "terragrunt-base", "modules", "terraform-aws-secrets-manager")
			return source, nil
		case "iam/policies":
			source := filepath.Join(commonDir, "terragrunt-base", "modules", "terraform-aws-iam/modules/iam-policy")
			return source, nil
		default:
			return "", fmt.Errorf("unsupported resource type: %s", resourceType)
		}
	}

	// If the starting directory does not contain a "common" directory, move up one level and try again
	parentDir := filepath.Dir(startingDir)
	if parentDir == startingDir {
		// We've reached the root directory and didn't find a "common" directory, so return an error
		return "", fmt.Errorf("could not find common directory")
	}

	return FindCommonDirectory(parentDir, resourceType)
}
