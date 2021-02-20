package gitcaddy

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Concurrency  int                      `yaml:"concurrency"`
	Repositories map[string][]*Repository `yaml:"repositories"`
	Identities   []*Identity              `yaml:"identities"`
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
