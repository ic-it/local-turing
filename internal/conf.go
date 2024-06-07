package internal

import (
	"io"

	"gopkg.in/yaml.v3"
)

const (
	DefaultURL = "https://www.turing.sk"
)

type (
	CloudTuring struct {
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
		URL      string `yaml:"url,omitempty"`
	}

	LocalTuringAssignment struct {
		Name          string   `yaml:"name"`
		Dir           string   `yaml:"dir,omitempty"`
		BuildCommands []string `yaml:"build-commands,omitempty"`
		Executable    string   `yaml:"executable,omitempty"`
		MainFile      string   `yaml:"main-file,omitempty"`
		PushName      string   `yaml:"push-name,omitempty"`
	}

	LocalTuring struct {
		TestsFile     string                  `yaml:"tests-file"`
		BuildCommands []string                `yaml:"build-commands,omitempty"`
		Executable    string                  `yaml:"executable,omitempty"`
		MainFile      string                  `yaml:"main-file,omitempty"`
		Assignments   []LocalTuringAssignment `yaml:"assignments"`
	}

	Config struct {
		CloudTuring CloudTuring `yaml:"cloud-turing,omitempty"`
		LocalTuring LocalTuring `yaml:"local-turing"`
	}
)

func ReadConfig(r io.Reader) (*Config, error) {
	var conf Config
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return nil, err
	}
	if err := conf.normalize(); err != nil {
		return nil, err
	}
	return &conf, nil
}

func (c *Config) normalize() error {
	if err := c.CloudTuring.normalize(); err != nil {
		return err
	}
	if err := c.LocalTuring.normalize(); err != nil {
		return err
	}
	return nil
}

func (ct *CloudTuring) normalize() error {
	if ct.URL == "" {
		ct.URL = DefaultURL
	}
	return nil
}

func (lt *LocalTuring) normalize() error {
	for i, lta := range lt.Assignments {
		if err := lta.normalize(lt.BuildCommands, lt.Executable, lt.MainFile); err != nil {
			return err
		}
		lt.Assignments[i] = lta
	}
	return nil
}

func (lta *LocalTuringAssignment) normalize(globalBuildCommands []string, globalExecutable, globalMainFile string) error {
	if lta.Dir == "" {
		lta.Dir = lta.Name
	}
	if len(lta.BuildCommands) == 0 {
		if len(globalBuildCommands) == 0 {
			return ErrNoBuildCommand
		}
		lta.BuildCommands = append(lta.BuildCommands, globalBuildCommands...)
	}
	if lta.Executable == "" {
		if globalExecutable == "" {
			return ErrNoExecutable
		}
		lta.Executable = globalExecutable
	}
	if lta.MainFile == "" {
		if globalMainFile == "" {
			return ErrNoMainFile
		}
		lta.MainFile = globalMainFile
	}
	if lta.PushName == "" {
		lta.PushName = lta.Name
	}
	return nil
}
