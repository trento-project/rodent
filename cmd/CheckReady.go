/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// CheckReadyCmd represents the CheckReady command
var CheckReadyCmd = &cobra.Command{
	Use:    "CheckReady",
	Short:  "Checks the status of the '/ready' endpoint",
	Long:   `The CheckReady command checks the status of the 'ready' endpoint and reports it to screen`,
	PreRun: DebugLogging,
	Run: func(cmd *cobra.Command, args []string) {
		CheckReadyFunc()
	},
}

func init() {
	rootCmd.AddCommand(CheckReadyCmd)
}
