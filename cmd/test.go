package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ic-it/local-turing/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:     "test [assignments]",
	Short:   "Test specific assignments or all assignments",
	Long:    `Test specific assignments or all assignments`,
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

		waitGroup := sync.WaitGroup{}
		writers := make([]*bufio.Writer, len(assignments))
		errors := make([]error, len(assignments))

		for i, assignment := range assignments {
			waitGroup.Add(1)
			logger.Debugw("running test", "assignment", assignment)
			tests, ok := tests.Tests[assignment.Name]
			if !ok {
				logger.Infow("no tests for assignment", "assignment", assignment)
				return ErrNoTests
			}
			go func(assignment internal.LocalTuringAssignment, tests []internal.Test, i int) {
				writers[i] = bufio.NewWriter(os.Stdout)
				errors[i] = processAssignment(assignment, tests, writers[i])
				waitGroup.Done()
			}(assignment, tests, i)
		}

		waitGroup.Wait()
		for _, w := range writers {
			w.Flush()
		}
		for _, err := range errors {
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func processAssignment(assignment internal.LocalTuringAssignment, tests []internal.Test, w *bufio.Writer) error {
	logger.Debugw("running tests", "assignment", assignment)
	logger.Debugw("build assignment", "assignment", assignment)
	if _, err := os.Stat(assignment.Dir); err != nil {
		logger.Infow("no such dir", "assignment", assignment)
		w.WriteString(color.RedString("No such dir for assignment %s\n", assignment.Name))
		return ErrNoSuchDir
	}
	w.WriteString(color.GreenString("\nBuilding assignment %s\n", assignment.Name))
	_, err := internal.RunBuild(context.Background(), assignment.BuildCommands, assignment.Dir, w, w)
	if err != nil {
		logger.Info("build failed")
		w.WriteString(color.RedString("Build failed for assignment %s\n", assignment.Name))
		return err
	}
	logger.Debugw("running tests", "assignment", assignment)
	testExecutable := filepath.Join(assignment.Dir, assignment.Executable)
	if stats, err := os.Stat(testExecutable); err != nil || !stats.Mode().IsRegular() {
		logger.Infow("no such executable", "assignment", assignment)
		w.WriteString(color.RedString("No such executable for assignment %s\n", assignment.Name))
		return ErrNoExecutable
	}
	for i, test := range tests {
		logger.Debugw("running test", "assignment", assignment, "test", test)
		sTime := time.Now()
		testRes, err := internal.RunTest(context.Background(), assignment.Dir, assignment.Executable, test.Inputs, test.Outputs)
		if err != nil {
			logger.Infow("test failed", "assignment", assignment, "test", test)
			w.WriteString(color.RedString("Test failed for assignment %s, test %d\n", assignment.Name, i))
			return err
		}
		if testRes.ExitCode != 0 {
			logger.Infow("test failed", "assignment", assignment, "test", test)
			w.WriteString(color.RedString("Test failed for assignment %s, test %d\n", assignment.Name, i))
			return ErrTestFailed
		}
		diffStr, err := internal.RenderResult(
			fmt.Sprintf("%s, test %d", assignment.Name, i),
			test.Outputs,
			testRes.StdOut,
			test.Inputs,
			err,
			time.Since(sTime),
			strip,
		)
		if err != nil {
			logger.Infow("test failed", "assignment", assignment, "test", test)
			w.WriteString(color.RedString("Test failed for assignment %s, test %d\n", assignment.Name, i))
			return err
		}
		w.WriteString(diffStr)
		w.WriteString("\n")
	}
	return nil
}

var (
	ErrNoTests          = errors.New("no tests for assignment")
	ErrNoSuchAssignment = errors.New("no such assignment")
	ErrTestFailed       = errors.New("test failed")
	ErrNoExecutable     = errors.New("no executable for assignment")
	ErrNoSuchDir        = errors.New("no such dir for assignment")
)
