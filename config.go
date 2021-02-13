package gitcaddy

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Concurrency  int                      `json:"concurrency"`
	Repositories map[string][]*Repository `json:"repositories"`
}

type Repository struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Enabled      *bool  `json:"enabled"`
	Destination  string `json:"destination"`
	Remote       string `json:"remote"`
	Depth        int    `json:"depth"`
	IdentityFile string `json:"identity_file"`
	AddFiles     string `json:"add_files"`
}

func (me *Repository) IsEnabled() bool {
	if me.Enabled == nil {
		return true
	}
	return *me.Enabled
}
func (me *Repository) Prefix() string {
	return fmt.Sprintf("[%s]", me.Name)
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

func (c *Config) LoadYaml(path string) error {
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

	if err := c.FixupConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) FixupConfig() error {

	if c.Concurrency == 0 {
		c.Concurrency = 5
	}
	return nil
}

func (c *Config) PrintConfig() error {
	buf, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", buf)
	return nil
}
