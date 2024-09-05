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
	log.Tracef("appended %d, total %d", len(ls), len(me.List))
}

func runRepositoryFile(opts *Options, configfile string, summary *RunSummary) error {
	cfg, err := LoadConfig(configfile)
	if err != nil {
		return err
	}

	run := &CompiledRun{}
	err = compileRepositoryFile(opts, configfile, summary, run)
	if err != nil {
		return err
	}

	// done compiling
	log.Debugf("done compiling repo list, have %d repositories to update", len(run.List))

	err = RunCompiled(opts, cfg, summary, run)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig(configfile string) (*gc.Config, error) {
	cfg := &gc.Config{}
	log.Infof("load config repository file:%s", configfile)
	err := cfg.LoadYaml(configfile)
	if err != nil {
		return nil, err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}
	return cfg, nil
}

// reads configfile and popultes the CompiledRun with repositories
// loads manifests as well
func compileRepositoryFile(opts *Options, configfile string, summary *RunSummary, run *CompiledRun) error {
	cfg, err := LoadConfig(configfile)
	if err != nil {
		return err
	}

	// run manifest if present
	if cfg.HasManifest() {
		err := CompileManifest(summary, opts, cfg, run)
		if err != nil {
			return err
		}
	}

	// compile repositories
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
