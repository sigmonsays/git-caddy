package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gc "github.com/sigmonsays/git-caddy"
	gologging "github.com/sigmonsays/go-logging"
)

type Options struct {
	LogLevel       string
	Section        string
	ConfigFile     string
	ManifestFile   string
	WorkingDir     string
	UpdateInterval int
	Action         string

	summary *RunSummary
}

func main() {
	opts := &Options{
		LogLevel:       "info",
		Section:        "",
		ConfigFile:     "repositories.yaml",
		ManifestFile:   filepath.Join(os.Getenv("HOME"), ".git-caddy.yaml"),
		WorkingDir:     "",
		UpdateInterval: 0,
		Action:         "config",
	}
	flag.StringVar(&opts.Action, "a", opts.Action, "specify processing actions (manfiest or config)")
	flag.StringVar(&opts.ConfigFile, "c", opts.ConfigFile, "specify config file")
	flag.StringVar(&opts.Section, "s", opts.Section, "section in config file")
	flag.StringVar(&opts.WorkingDir, "W", opts.WorkingDir, "initial working directory")
	flag.StringVar(&opts.LogLevel, "loglevel", opts.LogLevel, "log level")
	flag.StringVar(&opts.LogLevel, "l", opts.LogLevel, "short for -loglevel")
	flag.IntVar(&opts.UpdateInterval, "I", opts.UpdateInterval, "pull upstream for changes on an interval")
	flag.StringVar(&opts.ManifestFile, "m", opts.ManifestFile, "manifest file")
	flag.Parse()

	gologging.SetLogLevel(opts.LogLevel)

	var err error

	opts.summary = &RunSummary{}
	opts.summary.Start()

	manifest := &gc.ManifestConfig{}
	if gc.FileExists(opts.ManifestFile) {
		err = manifest.LoadYaml(opts.ManifestFile)
		ExitIfError(err, "LoadYaml %s: %s", opts.ManifestFile, err)
	}

	if opts.WorkingDir != "" {
		err = os.Chdir(opts.WorkingDir)
		ExitIfError(err, "Chdir %s: %s", opts.WorkingDir, err)
	}

	if opts.Action == "config" {
		if gc.FileExists(opts.ConfigFile) {
			err = runRepositoryFile(opts, opts.ConfigFile)
			ExitIfError(err, "run %s: %s", opts.ConfigFile, err)
		}
	}

	if opts.Action == "manifest" {
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
			err = runRepositoryFile(opts, e.Filename)
			if err != nil {
				log.Errorf("run %s: %s", e.Filename, err)
			}
		}
	}

	opts.summary.Stop()

	s := opts.summary
	log.Infof("scanned:%d errors:%d duration_sec:%d", s.Scanned, s.Errors, s.DurationSec)
}

func runRepositoryFile(opts *Options, configfile string) error {
	cfg := &gc.Config{}

	log.Infof("run repository file:%s section:%s", configfile, opts.Section)
	err := cfg.LoadYaml(configfile)
	if err != nil {
		return err
	}

	if log.IsTrace() {
		cfg.PrintConfig()
	}

	repos, found := cfg.Repositories[opts.Section]
	if found == false {
		return fmt.Errorf("Section not found: %q", opts.Section)
	}
	log.Debugf("concurrency:%d", cfg.Concurrency)

	updateRun := &UpdateRepositories{
		Section:      opts.Section,
		Cfg:          cfg,
		Repositories: repos,
		summary:      opts.summary,
	}

	if opts.UpdateInterval == 0 {
		err = updateRun.Run()
		return err
	}

	tick := time.NewTicker(time.Duration(opts.UpdateInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err = updateRun.Run()
			if err != nil {
				log.Warnf("%s", err)
			}
		}
	}

}
