package main

import (
	"os"

	gc "github.com/sigmonsays/git-caddy"
)

func RunManifest(summary *RunSummary, opts *Options, cfg *gc.Config) error {
	manifest := cfg.GetManifest()

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
		err := runRepositoryFile(opts, e.Filename, summary)
		if err != nil {
			log.Errorf("run %s: %s", e.Filename, err)
		}
	}
	return nil
}
