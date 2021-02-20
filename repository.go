package gitcaddy

import "fmt"

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

func (me *Repository) Validate() error {
	if me.Destination == "" {
		return fmt.Errorf("destination required")
	}
	if me.Remote == "" {
		return fmt.Errorf("remote required")
	}
	return nil
}
