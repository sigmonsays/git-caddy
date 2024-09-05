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

	// run manifest if present
	if cfg.HasManifest() {
		err := RunManifest(summary, opts, cfg)
		if err != nil {
			return err
		}
	}

	// run repositories
	compiled, err := CompileRepository(summary, opts, cfg)
	if err != nil {
		return err
	}

	err = compiled.Run()
	if err != nil {
		return err
	}

	return nil
}

func CompileRepository(summary *RunSummary, opts *Options, cfg *gc.Config) (*ProcessRepositories, error) {
	repos, found := cfg.Repositories[opts.Section]
	if found == false {
		return nil, fmt.Errorf("Section not found: %q", opts.Section)
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
		return nil, err
	}
	log.Tracef("compiled %d repos", len(compiled))
	processRun := &ProcessRepositories{}
	err = processRun.Run()
	if err != nil {
		return nil, err
	}
	return processRun, nil
}

func RunLoop(opts *Options, configfile string, summary *RunSummary) error {
	cfg := &gc.Config{}
	log.Infof("run repository file:%s section:%s", configfile, opts.Section)
	err := cfg.LoadYaml(configfile)
	if err != nil {
		return err
	}

	if log.IsTrace() {
		cfg.PrintConfig()
	}
	// run repositories
	compiled, err := CompileRepository(summary, opts, cfg)
	if err != nil {
		return err
	}
	return RunLoopCompiled(opts, compiled)
}

func RunLoopCompiled(opts *Options, compiled *ProcessRepositories) error {
	tick := time.NewTicker(time.Duration(opts.UpdateInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err := compiled.Run()
			if err != nil {
				log.Warnf("%s", err)
			}
		}
	}
	return nil
}
