package gitcaddy

import "fmt"

type Config struct {
	Repositories []*Repository `json:"repositories"`
}

type Repository struct {
	Destination  string `json:"destination"`
	Remote       string `json:"remote"`
	Depth        string `json:"depth"`
	IdentityFile string `json:"identity_file"`
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
