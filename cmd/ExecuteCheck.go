/*
Copyright © 2022 mr-stringer
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostToCheck, provider, checkID string
var callbackPort uint

// ExecuteCheckCmd represents the ExecuteCheck command
var ExecuteCheckCmd = &cobra.Command{
	Use:    "ExecuteCheck",
	Short:  "Executes a single check",
	Long:   `Executes a single check against a single host.`,
	PreRun: DebugLogging,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteCheckFunc()
	},
}

func init() {
	rootCmd.AddCommand(ExecuteCheckCmd)

	// Here you will define your flags and configuration settings.
	ExecuteCheckCmd.PersistentFlags().StringVarP(&hostToCheck, "hostToCheck", "t", "", "The host that the runner will execute the check on")
	ExecuteCheckCmd.PersistentFlags().StringVarP(&checkID, "checkID", "c", "", "The ID of the check to be run")
	ExecuteCheckCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "default", "The provider the the host runs on. default, aws, gcp and azure are the only supported values")
	ExecuteCheckCmd.PersistentFlags().UintVar(&callbackPort, "callbackPort", 8000, "The port that the callback listener will use, the default value is 8000")

	err := ExecuteCheckCmd.MarkPersistentFlagRequired("hostToCheck")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = ExecuteCheckCmd.MarkPersistentFlagRequired("checkID")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ExecuteCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ExecuteCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
