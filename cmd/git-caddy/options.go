package main

import (
	"os"

	gologging "github.com/sigmonsays/go-logging"
	"github.com/spf13/cobra"
)

func DefaultOptions() *Options {
	opts := &Options{
		ConfigFile:     "repositories.yaml",
		WorkingDir:     "",
		UpdateInterval: 0,
		Section:        "repos",
		LogLevel:       "info",
	}
	return opts
}

func ReadOptions(cmd *cobra.Command) (*Options, error) {
	opts := DefaultOptions()

	opts.Section, _ = cmd.Flags().GetString("section")
	opts.ConfigFile, _ = cmd.Flags().GetString("config")
	opts.WorkingDir, _ = cmd.Flags().GetString("workdir")
	opts.LogLevel, _ = cmd.Flags().GetString("loglevel")
	opts.UpdateInterval, _ = cmd.Flags().GetInt("interval")
	opts.Pretend, _ = cmd.Flags().GetBool("pretend")

	if opts.LogLevel != "" {
		gologging.SetLogLevel(opts.LogLevel)
	}

	if opts.WorkingDir != "" {
		err := os.Chdir(opts.WorkingDir)
		ExitIfError(err, "Chdir %s: %s", opts.WorkingDir, err)
		log.Debugf("changed working directory to %s", opts.WorkingDir)
	}

	return opts, nil
}

type Options struct {
	LogLevel       string
	Section        string
	ConfigFile     string
	WorkingDir     string
	UpdateInterval int
	Discover       bool
	Version        bool
	Pretend        bool
}
