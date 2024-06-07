package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/ic-it/local-turing/internal"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	// force colors to be used
	forceColors bool
	// configFile is the path to the config file
	configFile string = "config.yml"
	// config
	config *internal.Config
	// verbose is the flag to print the output of the test cases
	verbose bool
	// strip is the flag to strip the output and expected output before comparing
	strip bool
	// logger is the logger (lol)
	logger *zap.SugaredLogger
	// tests is the tests (`\_(ツ)_/¯)
	tests internal.Tests
)

var (
	_ = config // FIXME: remove this line
	_ = logger // FIXME: remove this line
	_ = tests  // FIXME: remove this line
)

func preRunE(*cobra.Command, []string) (err error) {
	if forceColors {
		color.NoColor = false
		// color.Output = os.Stdout
	}
	internal.SetupLogger(verbose)
	logger = internal.GetLogger()
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	config, err = internal.ReadConfig(file)
	if err != nil {
		logger.Infow("error reading config", "error", err)
		return err
	}
	logger.Debug("loaded config")

	testsFile, err := os.Open(config.LocalTuring.TestsFile)
	if err != nil {
		return err
	}
	defer testsFile.Close()
	tests, err = internal.TestsUnmarshal(testsFile)
	if err != nil {
		logger.Error("error reading tests file", err)
		return err
	}
	logger.Debug("loaded tests")

	logger.Debugw(
		"setup",
		"forceColors", forceColors,
		"verbose", verbose,
		"strip", strip,
	)
	return nil
}
