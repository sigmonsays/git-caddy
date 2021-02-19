package gitcaddy

import "fmt"

type Identity struct {
	Name         string   `yaml:"name"`
	FullName     string   `yaml:"full_name"`
	Email        string   `yaml:"email"`
	Repositories []string `yaml:"repositories"`
}

func (me *Identity) Validate() error {
	if me.Name == "" {
		return fmt.Errorf("name required")
	}
	if me.FullName == "" {
		return fmt.Errorf("full_name required")
	}
	if me.Email == "" {
		return fmt.Errorf("email required")
	}
	return nil
}
