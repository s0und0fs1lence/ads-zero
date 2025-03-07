package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ads-zero",
	Short: "Ads-Zero is a tool for fetching spend metrics from Advertising platform, and send alerts based on rules",

	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here
		if len(args) < 1 {
			return fmt.Errorf("provide")
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
