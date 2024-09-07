package main

import (
	"time"

	gc "github.com/sigmonsays/git-caddy"
)

func RunLoop(opts *Options, configfile string, summary *RunSummary) error {
	cfg := &gc.Config{}
	log.Infof("run loop config:%s section:%s", configfile, opts.Section)
	log.Tracef("opts %+v", opts)
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
			log.Infof("running ... ")
			err := RunCompiled(opts, cfg, summary, run)
			if err != nil {
				log.Warnf("%s", err)
			}
		}
	}
	return nil
}
