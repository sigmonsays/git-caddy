package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	gc "github.com/sigmonsays/git-caddy"
)

// These variables are populated by goreleaser when the binary is built.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "git-caddy",
		Short: "manage git repositories",
	}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "run config file",
		Run: func(cmd *cobra.Command, args []string) {
			summary := NewRunSummary()
			summary.Start()
			defer FinishSummary(summary)
			opts, err := ReadOptions(cmd)
			if err != nil {
				ExitError("ReadOptions: %s", err)
			}
			if gc.FileExists(opts.ConfigFile) {
				err := runRepositoryFile(opts, opts.ConfigFile, summary)
				ExitIfError(err, "run %s: %s", opts.ConfigFile, err)
			}
		},
	}
	rootCmd.AddCommand(runCmd)

	var discoverCmd = &cobra.Command{
		Use:   "discover",
		Short: "discover repositories",
		Run: func(cmd *cobra.Command, args []string) {
			opts, err := ReadOptions(cmd)
			if err != nil {
				ExitError("ReadOptions: %s", err)
			}
			err = doDiscovery(opts.WorkingDir)
			if err != nil {
				ExitError("Discovery: %s", err)
			}
		},
	}
	rootCmd.AddCommand(discoverCmd)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version %s\n", version)
			fmt.Printf("commit %s\n", commit)
			fmt.Printf("date %s\n", date)
			os.Exit(0)
		},
	}
	rootCmd.AddCommand(versionCmd)

	var manifestCmd = &cobra.Command{
		Use:   "manifest",
		Short: "run manifest",
		Run: func(cmd *cobra.Command, args []string) {
			summary := NewRunSummary()
			summary.Start()
			defer FinishSummary(summary)
			opts, err := ReadOptions(cmd)
			if err != nil {
				ExitError("ReadOptions: %s", err)
			}

			manifest := &gc.ManifestConfig{}
			if gc.FileExists(opts.ManifestFile) {
				err = manifest.LoadYaml(opts.ManifestFile)
				ExitIfError(err, "LoadYaml %s: %s", opts.ManifestFile, err)
			}

			files := manifest.ListManifest()
			log.Tracef("loaded %d files using manifest from %s", len(files), opts.ManifestFile)
			for _, e := range files {
				if e.Section != "" {
					opts.Section = e.Section
				}
				if e.Def.WorkingDir != "" {
					workingdir := os.ExpandEnv(e.Def.WorkingDir)
					log.Tracef("chdir %s", workingdir)
					os.Chdir(workingdir)
				}
				err = runRepositoryFile(opts, e.Filename, summary)
				if err != nil {
					log.Errorf("run %s: %s", e.Filename, err)
				}
			}
		},
	}
	rootCmd.AddCommand(manifestCmd)

	dopts := DefaultOptions()
	rootCmd.PersistentFlags().StringP("loglevel", "l", dopts.LogLevel, "log level")
	rootCmd.PersistentFlags().StringP("section", "s", dopts.Section, "section name")
	rootCmd.PersistentFlags().StringP("config", "c", dopts.ConfigFile, "repositories.yaml")
	rootCmd.PersistentFlags().StringP("workdir", "W", dopts.WorkingDir, "working directory")
	rootCmd.PersistentFlags().Int32P("interval", "i", int32(dopts.UpdateInterval), "update interval")
	rootCmd.PersistentFlags().StringP("manifest", "m", dopts.ManifestFile, "manifest file")

	err := rootCmd.Execute()
	if err != nil {
		ExitError("%s", err)
	}
}
