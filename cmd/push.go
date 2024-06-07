package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:     "push [assignments]",
	Short:   "Push specific assignments or all assignments",
	Long:    `Push specific assignments or all assignments`,
	PreRunE: preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented yet")
	},
}
