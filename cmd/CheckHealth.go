/*
Copyright © 2022 mr-stringer
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// CheckHealthCmd represents the CheckHealth command
var CheckHealthCmd = &cobra.Command{
	Use:    "CheckHealth",
	Short:  "Checks the status of the '/health' endpoint",
	Long:   `The CheckReady command checks the status of the 'health' endpoint and reports it to screen`,
	PreRun: DebugLogging,
	Run: func(cmd *cobra.Command, args []string) {
		CheckHealthFunc()
	},
}

func init() {
	rootCmd.AddCommand(CheckHealthCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// CheckHealthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// CheckHealthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
