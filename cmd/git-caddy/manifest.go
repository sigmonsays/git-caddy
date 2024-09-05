package main

import (
	"os"

	gc "github.com/sigmonsays/git-caddy"
)

func CompileManifest(summary *RunSummary, opts *Options, cfg *gc.Config, run *CompiledRun) error {
	manifest := cfg.GetManifest()

	files := manifest.ListManifest()
	log.Tracef("loaded %d files using manifest", len(files))
	for _, e := range files {
		if e.Section != "" {
			opts.Section = e.Section
		}
		if e.Def.WorkingDir != "" {
			workingdir := os.ExpandEnv(e.Def.WorkingDir)
			log.Tracef("chdir %s", workingdir)
			os.Chdir(workingdir)
		}
		err := CompileRepository(summary, opts, cfg, run)
		if err != nil {
			log.Errorf("run %s: %s", e.Filename, err)
		}
	}
	return nil
}
