package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "local-turing",
	Short:   "A simple CLI tool to run your tests on your local machine",
	Long:    `A simple CLI tool to run your tests on your local machine`,
	Version: "0.1.5-alpha",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&forceColors, "force-colors", "C", forceColors, "force colors to be used")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "print the output of the test cases")
	rootCmd.PersistentFlags().BoolVarP(&strip, "strip", "s", strip, "strip the output and expected output before comparing")
}
