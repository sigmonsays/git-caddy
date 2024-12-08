package main

import (
	"fmt"
	"path/filepath"

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

// main entry point from cobra for 'git-caddy run'
func runRepositoryFile(opts *Options, configfile string, summary *RunSummary) error {
	log.Infof("Load config %s", configfile)
	cfg, err := LoadConfig(configfile)
	if err != nil {
		return err
	}

	abspath, err := filepath.Abs(configfile)
	if err != nil {
		return err
	}
	if opts.Dir == "" {
		opts.Dir = filepath.Dir(abspath)
		log.Debug("repofile %s: Defaulting dir to %s",
			configfile, opts.Dir)
	}
	log.Infof("runRepositoryFile %s (dir:%s)", configfile, opts.Dir)

	run := &CompiledRun{}
	err = compileRepositoryFile(opts, cfg, configfile, summary, run)
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

// reads configfile and popultes the CompiledRun with repositories
// loads manifests as well
func compileRepositoryFile(opts *Options, cfg *gc.Config, configfile string, summary *RunSummary, run *CompiledRun) error {

	// run manifest if present
	if cfg.HasManifest() {
		log.Tracef("has %d manifest entries", len(cfg.Manifest))
		err := CompileManifest(summary, opts, cfg, run)
		if err != nil {
			return err
		}
	}

	// compile repositories
	if len(cfg.Repositories) > 0 {
		err := CompileRepository(summary, opts, cfg, run)
		if err != nil {
			return err
		}
	}

	return nil
}

func CompileRepository(summary *RunSummary, opts *Options, cfg *gc.Config, run *CompiledRun) error {
	log.Tracef("compiling repository section %s", opts.Section)
	repos, found := cfg.Repositories[opts.Section]
	if found == false {
		return fmt.Errorf("Section not found: %q", opts.Section)
	}
	log.Debugf("concurrency:%d", cfg.Concurrency)

	compile := &Compile{
		Section:      opts.Section,
		Cfg:          cfg,
		Dir:          opts.Dir,
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

func LoadConfig(configfile string) (*gc.Config, error) {
	cfg := &gc.Config{}
	log.Debugf("load config %s", configfile)
	err := cfg.LoadYaml(configfile)
	if err != nil {
		return nil, err
	}
	// if log.IsTrace() {
	// 	cfg.PrintConfig()
	// }
	return cfg, nil
}
