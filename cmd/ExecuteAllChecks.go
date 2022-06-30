/*
Copyright © 2022 mr-stringer
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var checkInterval int

// ExecuteAllChecksCmd represents the ExecuteAllChecks command
var ExecuteAllChecksCmd = &cobra.Command{
	Use:    "ExecuteAllChecks",
	Short:  "Executes all checks",
	Long:   `Executes all the checks found in the /catalog endpoint against a single target`,
	PreRun: DebugLogging,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteAllChecksFunc()
	},
}

func init() {
	rootCmd.AddCommand(ExecuteAllChecksCmd)
	/*These flags already exist in 'ExecuteCheck' */
	ExecuteAllChecksCmd.PersistentFlags().StringVarP(&runner, "runner", "r", "", "trento runner to test (IP address or hostname)")
	ExecuteAllChecksCmd.PersistentFlags().StringVarP(&hostToCheck, "hostToCheck", "t", "", "The host that the runner will execute the check on")
	ExecuteAllChecksCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "default", "The provider the the host runs on. default, aws, gcp and azure are the only supported values, default is the default")
	ExecuteAllChecksCmd.PersistentFlags().IntVarP(&checkInterval, "checkInterval", "i", 5, "The pause between submitting check requests in seconds, default 5")
	ExecuteAllChecksCmd.PersistentFlags().StringVarP(&callbackUrl, "callbackUrl", "u", "/api/runner/callbacks", "The url to listen to use for callbacks, only required if the runner is using a custom url")
	ExecuteAllChecksCmd.PersistentFlags().UintVar(&callbackPort, "callbackPort", 8000, "The port that the callback listener will use, the default value is 8000")

	err := ExecuteAllChecksCmd.MarkPersistentFlagRequired("hostToCheck")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ExecuteAllChecksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ExecuteAllChecksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
