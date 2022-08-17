package gitcaddy

import (
	"fmt"
	"path/filepath"
	"strings"

	giturl "github.com/whilp/git-urls"
)

type Repository struct {
	Section      string   `yaml:"section"`
	Name         string   `yaml:"name"`
	Names        []string `yaml:"names"`
	Description  string   `yaml:"description"`
	Enabled      *bool    `yaml:"enabled"`
	Destination  string   `yaml:"destination"`
	Remote       string   `yaml:"remote"`
	Depth        int      `yaml:"depth"`
	IdentityFile string   `yaml:"identity_file"`
	AddFiles     string   `yaml:"add_files"`
	NoClone      bool     `yaml:"no_clone"`
}

func (me *Repository) Copy() *Repository {
	cp := &Repository{}
	cp.Section = me.Section
	cp.Name = me.Name
	cp.Names = me.Names
	cp.Description = me.Description
	cp.Enabled = me.Enabled
	cp.Destination = me.Destination
	cp.Remote = me.Remote
	cp.Depth = me.Depth
	cp.IdentityFile = me.IdentityFile
	cp.AddFiles = me.AddFiles
	cp.NoClone = me.NoClone
	return cp
}
func (me *Repository) IsEnabled() bool {
	if me.Enabled == nil {
		return true
	}
	return *me.Enabled
}

func (me *Repository) Prefix(sub string) string {
	return fmt.Sprintf("[%s %s] ", me.Name, sub)
}

func (me *Repository) Defaults() error {

	// fill in name and destination from the remote if possible
	if me.Name == "" || me.Destination == "" {
		// // to parse as a url we need a prefix://
		// i := strings.Index(me.Remote, "//:")
		// var remote string
		// if i == -1 {
		// 	remote = "default://" + me.Remote
		// } else {
		// 	remote = me.Remote
		// }
		p, err := giturl.Parse(me.Remote)
		if err != nil {
			return err
		}
		basename := filepath.Base(p.Path)
		if strings.HasSuffix(basename, ".git") {
			basename = strings.TrimSuffix(basename, ".git")
		}
		me.Name = basename
		me.Destination = "./" + basename

		return nil
	}
	return nil
}

func (me *Repository) Validate() error {
	if me.Destination == "" {
		return fmt.Errorf("destination required")
	}
	if me.Remote == "" {
		return fmt.Errorf("remote required")
	}
	return nil
}
