package main

import (
	"os"

	gc "github.com/sigmonsays/git-caddy"
)

func CompileManifest(summary *RunSummary, opts *Options, cfg *gc.Config, run *CompiledRun) error {
	manifest := cfg.GetManifest()

	entries := manifest.ListManifest()
	log.Tracef("loaded %d manifest entries from manifest", len(entries))
	for _, e := range entries {
		if e.Section != "" {
			opts.Section = e.Section
		}
		if e.Def.WorkingDir != "" {
			workingdir := os.ExpandEnv(e.Def.WorkingDir)
			log.Tracef("chdir %s", workingdir)
			os.Chdir(workingdir)
		}

		cfg2, err := LoadConfig(e.Filename)
		if err != nil {
			log.Errorf("load %s: %s", e.Filename, err)
			continue
		}
		opts2 := opts.Clone()
		opts2.Section = e.Section

		err = CompileRepository(summary, opts2, cfg2, run)
		if err != nil {
			log.Errorf("run %s: %s", e.Filename, err)
		}
	}
	return nil
}
