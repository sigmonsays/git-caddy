package gitcaddy

import (
	"fmt"
	"path/filepath"
	"strings"

	giturl "github.com/whilp/git-urls"
)

type Repository struct {
	Section      string `yaml:"section"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Enabled      *bool  `yaml:"enabled"`
	Destination  string `yaml:"destination"`
	Remote       string `yaml:"remote"`
	Depth        int    `yaml:"depth"`
	IdentityFile string `yaml:"identity_file"`
	AddFiles     string `yaml:"add_files"`
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
