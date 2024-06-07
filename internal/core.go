package internal

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type TestResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func RunTest(ctx context.Context, workingDir, executable, testInput, testOutput string) (*TestResult, error) {
	logger.Debug("running test", "executable", executable, "workingDir", workingDir)
	cmd := exec.CommandContext(ctx, executable, testInput)
	var stdout, stderr strings.Builder
	cmd.Stdin = strings.NewReader(testInput)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = workingDir
	if err := cmd.Run(); err != nil {
		logger.Error("error running test", "error", err)
		return nil, err
	}
	if cmd.ProcessState.ExitCode() != 0 {
		logger.Error("test failed by exit code", "exitCode", cmd.ProcessState.ExitCode(), "stdout", stdout.String(), "stderr", stderr.String())
		return &TestResult{StdOut: stdout.String(), StdErr: stderr.String(), ExitCode: cmd.ProcessState.ExitCode()}, ErrTestFailedByExitCode
	}
	return &TestResult{
		StdOut:   stdout.String(),
		StdErr:   stderr.String(),
		ExitCode: cmd.ProcessState.ExitCode(),
	}, nil
}

type BuildResult struct {
	ExitCode int
}

func RunBuild(ctx context.Context, buildCommands []string, workingDir string, stdOutWriter io.Writer, stdErrWriter io.Writer) (*BuildResult, error) {
	logger.Debug("running build", "buildCommands", buildCommands)
	for _, buildCommand := range buildCommands {
		fmt.Fprintf(stdOutWriter, "> %s\n", buildCommand)
		args := strings.Split(buildCommand, " ")
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Stdout = stdOutWriter
		cmd.Stderr = stdErrWriter
		cmd.Dir = workingDir
		if err := cmd.Run(); err != nil {
			logger.Infow("error running build", "error", err)
			return nil, err
		}
		if cmd.ProcessState.ExitCode() != 0 {
			logger.Infow("build failed", "exitCode", cmd.ProcessState.ExitCode())
			return &BuildResult{ExitCode: cmd.ProcessState.ExitCode()}, ErrBuildFailed
		}
	}
	return &BuildResult{ExitCode: 0}, nil
}
