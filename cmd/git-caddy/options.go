package main

import (
	gologging "github.com/sigmonsays/go-logging"
	"github.com/spf13/cobra"
)

func DefaultOptions() *Options {
	opts := &Options{
		ConfigFile:     "repositories.yaml",
		WorkingDir:     "",
		UpdateInterval: 60,
		Section:        "repos",
		LogLevel:       "info",
	}
	return opts
}

func ReadOptions(cmd *cobra.Command, args []string) (*Options, error) {
	opts := DefaultOptions()

	opts.Section, _ = cmd.Flags().GetString("section")
	opts.WorkingDir, _ = cmd.Flags().GetString("workdir")
	opts.LogLevel, _ = cmd.Flags().GetString("loglevel")
	opts.UpdateInterval, _ = cmd.Flags().GetInt32("interval")
	opts.Pretend, _ = cmd.Flags().GetBool("pretend")

	if opts.LogLevel != "" {
		gologging.SetLogLevel(opts.LogLevel)
	}
	// first argument is config file
	if len(args) > 0 {
		opts.ConfigFile = args[0]
	} else {
		opts.ConfigFile = "repositories.yaml"
	}

	return opts, nil
}

type Options struct {
	LogLevel       string
	Section        string
	ConfigFile     string
	WorkingDir     string
	UpdateInterval int32
	Discover       bool
	Version        bool
	Pretend        bool
}

func (me *Options) Clone() *Options {
	var opts2 Options
	opts2 = *me
	return &opts2
}
