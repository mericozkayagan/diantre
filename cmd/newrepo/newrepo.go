/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package newrepo

import (
	"fmt"

	"github.com/mericozkayagan/diantre/src/newrepo"
	"github.com/spf13/cobra"
)

// newrepoCmd represents the newrepo command
var NewrepoCmd = &cobra.Command{
	Use:   "newrepo",
	Short: "Command needed to create a new repo",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("newrepo called")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newrepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newrepoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newrepo.CreateNewRepo();
}
