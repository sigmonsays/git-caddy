package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func doDiscovery(directory string) error {

	// empty string means current directory
	if directory == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		directory = wd
	}
	names, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	repos := make([]*repository, 0)
	for _, name := range names {
		if name.IsDir() == false {
			continue
		}
		basename := name.Name()
		repodir := filepath.Join(directory, basename)
		tpath := filepath.Join(repodir, ".git")
		st, err := os.Stat(tpath)
		if err != nil {
			continue
		}
		if st.IsDir() == false {
			continue
		}

		// git remote get-url --push origin
		buf := bytes.NewBuffer(nil)
		c := exec.Command("git", "remote", "get-url", "--push", "origin")
		c.Dir = repodir
		c.Stdout = buf
		err = c.Run()
		if err != nil {
			continue
		}
		remote := strings.Trim(buf.String(), "\n")
		if remote == "" {
			continue
		}
		log.Debugf("path %s has remote %s", repodir, remote)
		repo := &repository{}
		repo.Name = basename
		repo.Remote = remote
		repos = append(repos, repo)
	}
	cfg := &config{}
	cfg.Repositories = make(map[string][]*repository, 0)
	cfg.Repositories["discovered"] = repos
	buf, _ := yaml.Marshal(cfg)
	fmt.Printf("%s\n", buf)

	return nil
}

type config struct {
	Repositories map[string][]*repository
}

type repository struct {
	Name   string `yaml:"name"`
	Remote string `yaml:"remote"`
}
