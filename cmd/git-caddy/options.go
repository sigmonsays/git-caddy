package main

import (
	"os"
	"path/filepath"

	gologging "github.com/sigmonsays/go-logging"
	"github.com/spf13/cobra"
)

func DefaultOptions() *Options {
	opts := &Options{
		ConfigFile:     "repositories.yaml",
		ManifestFile:   filepath.Join(os.Getenv("HOME"), ".git-caddy.yaml"),
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
	opts.ManifestFile, _ = cmd.Flags().GetString("manifest")

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
	ManifestFile   string
	WorkingDir     string
	UpdateInterval int
	Discover       bool
	Version        bool
}
