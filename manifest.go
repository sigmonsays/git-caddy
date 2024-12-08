package gitcaddy

import (
	"os"
	"path/filepath"
	"strings"
)

type ManifestDef struct {
	Pattern    string `yaml:"pattern"`
	Sections   string `yaml:"sections"`
	WorkingDir string `yaml:"dir"`
}

// after expanding the glob pattern, we yield a ManifestEntry
type ManifestEntry struct {
	Def      *ManifestDef
	Filename string
	Section  string
}

type ManifestConfig struct {
	Manifest []*ManifestDef
}

func (c *ManifestConfig) ListManifest() []*ManifestEntry {
	var ret []*ManifestEntry
	for _, e := range c.Manifest {
		pattern := os.ExpandEnv(e.Pattern)
		matches, err := filepath.Glob(pattern)
		log.Tracef("glob %s got %d matches", e.Pattern, len(matches))
		if err != nil {
			log.Warnf("Glob %s: %s", e.Pattern, err)
			continue
		}
		for _, match := range matches {
			sections := strings.Fields(e.Sections)
			for _, section := range sections {
				if e.WorkingDir == "" {
					abs, err := filepath.Abs(match)
					if err == nil {
						e.WorkingDir = filepath.Dir(abs)
						log.Debugf("manifest entry %s: Setting default workingdir to %s",
							match, e.WorkingDir)
					}
				}
				ent := &ManifestEntry{}
				ent.Filename = match
				ent.Section = section
				ent.Def = e
				ret = append(ret, ent)
			}
		}
	}
	return ret
}
