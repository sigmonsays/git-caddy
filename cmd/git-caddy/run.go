package main

import (
	"fmt"
	"time"

	gc "github.com/sigmonsays/git-caddy"
)

// a complete run of all repositories once manifest and configuration files have been parsed
type CompiledRun struct {
	List []*CompiledRepository
}

func (me *CompiledRun) Append(ls []*CompiledRepository) {
	me.List = append(me.List, ls...)
}

func runRepositoryFile(opts *Options, configfile string, summary *RunSummary) error {
	run := &CompiledRun{}
	err := compileRepositoryFile(opts, configfile, summary, run)
	if err != nil {
		return err
	}
	return nil
}

func compileRepositoryFile(opts *Options, configfile string, summary *RunSummary, run *CompiledRun) error {
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
		err := CompileManifest(summary, opts, cfg, run)
		if err != nil {
			return err
		}
	}

	// run repositories
	err = CompileRepository(summary, opts, cfg, run)
	if err != nil {
		return err
	}

	return nil
}

func CompileRepository(summary *RunSummary, opts *Options, cfg *gc.Config, run *CompiledRun) error {
	repos, found := cfg.Repositories[opts.Section]
	if found == false {
		return fmt.Errorf("Section not found: %q", opts.Section)
	}
	log.Debugf("concurrency:%d", cfg.Concurrency)

	compile := &Compile{
		Section:      opts.Section,
		Cfg:          cfg,
		Repositories: repos,
		summary:      summary,
	}

	compiled, err := compile.Run()
	if err != nil {
		return err
	}
	run.Append(compiled)
	log.Tracef("compiled %d repos", len(compiled))
	return nil
}

func RunLoop(opts *Options, configfile string, summary *RunSummary) error {
	cfg := &gc.Config{}
	log.Infof("run repository file:%s section:%s", configfile, opts.Section)
	err := cfg.LoadYaml(configfile)
	if err != nil {
		return err
	}
	run := &CompiledRun{}

	if log.IsTrace() {
		cfg.PrintConfig()
	}
	// run repositories
	err = CompileRepository(summary, opts, cfg, run)
	if err != nil {
		return err
	}
	return RunLoopCompiled(cfg, summary, opts, run)
}

func RunLoopCompiled(cfg *gc.Config, summary *RunSummary, opts *Options, run *CompiledRun) error {
	tick := time.NewTicker(time.Duration(opts.UpdateInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err := RunCompiled(opts, cfg, summary, run)
			if err != nil {
				log.Warnf("%s", err)
			}
		}
	}
	return nil
}
