/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// CheckCatalogCmd represents the CheckCatalog command
var CheckCatalogCmd = &cobra.Command{
	Use:    "CheckCatalog",
	Short:  "Retrieves all check ID found on the runner",
	Long:   `Retrieves all check ID found on the runner using the /catalog endpoint.  These IDs can then be used with the "ExecuteCheck" command.`,
	PreRun: DebugLogging,
	Run: func(cmd *cobra.Command, args []string) {
		CheckCatalogFunc()
	},
}

func init() {
	rootCmd.AddCommand(CheckCatalogCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// CheckCatalogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// CheckCatalogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
