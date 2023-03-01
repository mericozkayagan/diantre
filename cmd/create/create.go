/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"
	"strings"

	"github.com/mericozkayagan/diantre/src/create"
	"github.com/spf13/cobra"
)

var (
	createType string
)

// deployCmd represents the deploy command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify what you will create: S3, Secrets-Manager")
		fmt.Scanln(&createType)

		switch strings.ToLower(createType) {
		case "s3":
			create.CreateS3()
		case "secrets-manager":
			create.CreateSecretsManager()
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
