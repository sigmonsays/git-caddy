package main

import (
	"fmt"
	"time"

	gc "github.com/sigmonsays/git-caddy"
)

func runRepositoryFile(opts *Options, configfile string, summary *RunSummary) error {
	cfg := &gc.Config{}

	log.Infof("run repository file:%s section:%s", configfile, opts.Section)
	err := cfg.LoadYaml(configfile)
	if err != nil {
		return err
	}

	if log.IsTrace() {
		cfg.PrintConfig()
	}

	hasManifest := cfg.HasManifest()
	if hasManifest {
		return RunManifest(summary, opts, cfg)
	}
	return RunRepository(summary, opts, cfg)
}

func RunRepository(summary *RunSummary, opts *Options, cfg *gc.Config) error {
	repos, found := cfg.Repositories[opts.Section]
	if found == false {
		return fmt.Errorf("Section not found: %q", opts.Section)
	}
	log.Debugf("concurrency:%d", cfg.Concurrency)

	compile := &CompiledRun{
		Section:      opts.Section,
		Cfg:          cfg,
		Repositories: repos,
		summary:      summary,
	}

	compiled, err := compile.Run()
	if err != nil {
		return err
	}
	log.Tracef("compiled %d repos", len(compiled))

	processRun := &ProcessRepositories{}

	if opts.UpdateInterval == 0 {
		err := processRun.Run(compiled)
		return err
	}

	// run loop
	tick := time.NewTicker(time.Duration(opts.UpdateInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err := processRun.Run(compiled)
			if err != nil {
				log.Warnf("%s", err)
			}
		}
	}

	return nil
}
