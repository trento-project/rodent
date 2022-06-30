/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "development"

// VersionCmd represents the Version command
var VersionCmd = &cobra.Command{
	Use:   "Version",
	Short: "Display the version of rodent",
	Long:  "Display the version of rodent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(VersionCmd)
}
