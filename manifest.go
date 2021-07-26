package gitcaddy

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type ManifestConfig struct {
	RepositoryFiles []*ManifestDef `yaml:"repository_files"`
}
type ManifestDef struct {
	Pattern  string
	Sections string
}
type ManifestEntry struct {
	Filename string
	Section  string
}

func (c *ManifestConfig) LoadYaml(path string) error {
	log.Tracef("load yaml file %s", path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}

	if err := c.LoadYamlBuffer(b.Bytes()); err != nil {
		return err
	}

	if err := c.Fixup(); err != nil {
		return err
	}

	return nil
}

func (c *ManifestConfig) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ManifestConfig) Fixup() error {
	return nil
}

func (c *ManifestConfig) Print() error {
	buf, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", buf)
	return nil
}

func (c *ManifestConfig) ListManifest() []*ManifestEntry {
	var ret []*ManifestEntry
	for _, e := range c.RepositoryFiles {
		matches, err := filepath.Glob(e.Pattern)
		log.Tracef("glob %s got %d matches", e.Pattern, len(matches))
		if err != nil {
			log.Warnf("Glob %s: %s", e.Pattern, err)
			continue
		}
		for _, match := range matches {
			sections := strings.Fields(e.Sections)
			for _, section := range sections {
				e := &ManifestEntry{}
				e.Filename = match
				e.Section = section
				ret = append(ret, e)
			}
		}
	}
	return ret
}
