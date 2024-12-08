package main

import (
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

		cfg2, err := LoadConfig(e.Filename)
		if err != nil {
			log.Errorf("load %s: %s", e.Filename, err)
			continue
		}
		opts2 := opts.Clone()
		opts2.Section = e.Section
		opts2.Dir = e.Def.WorkingDir

		// merge identities together
		cfg2.Identities = append(cfg2.Identities, cfg.Identities...)

		err = CompileRepository(summary, opts2, cfg2, run)
		if err != nil {
			log.Errorf("run %s: %s", e.Filename, err)
		}
	}
	return nil
}
