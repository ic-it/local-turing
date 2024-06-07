package internal

import "errors"

var (
	ErrTestFailedByExitCode = errors.New("test failed by exit code")
	ErrBuildFailed          = errors.New("build failed")
	ErrNoBuildCommand       = errors.New("no build command specified")
	ErrNoExecutable         = errors.New("no executable specified")
	ErrNoMainFile           = errors.New("no main file specified")
	ErrInvalidTests         = errors.New("invalid tests")
)
