package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ic-it/local-turing/internal"
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
		var assignments []internal.LocalTuringAssignment
		if len(args) > 0 {
			for _, arg := range args {
				found := false
				for _, assignment := range config.LocalTuring.Assignments {
					if assignment.Name == arg {
						assignments = append(assignments, assignment)
						found = true
					}
				}
				if !found {
					logger.Infow("no such assignment", "assignment", arg)
					return ErrNoSuchAssignment
				}
			}
		} else {
			assignments = config.LocalTuring.Assignments
		}
		client, err := internal.NewTuringClient(internal.DefaultURL)
		if err != nil {
			return err
		}
		err = client.Login(config.CloudTuring.Name, config.CloudTuring.Password)
		if err != nil {
			return err
		}
		for _, assignment := range assignments {
			fmt.Printf("Pushing %s", assignment.Name)
			err = client.SaveAssigment(&assignment)
			if err != nil {
				logger.Errorw("push failed", "assignment", assignment, "error", err)
				return err
			}
			logger.Infow("push succeeded", "assignment", assignment)
			color.Green("\t . . . Push succeeded")
			fmt.Println("Assignment link:", client.GetAssignmentLink(&assignment))
		}
		return nil
	},
}
