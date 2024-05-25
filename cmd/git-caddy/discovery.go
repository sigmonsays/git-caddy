package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gc "github.com/sigmonsays/git-caddy"
)

func doDiscovery() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	names, err := os.ReadDir(wd)
	if err != nil {
		return err
	}
	repos := make([]*gc.Repository, 0)
	for _, name := range names {
		if name.IsDir() == false {
			continue
		}
		basename := name.Name()
		repodir := filepath.Join(wd, basename)
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
		repo := &gc.Repository{}
		repo.Name = basename
		repo.Remote = remote
		repos = append(repos, repo)
	}
	cfg := &gc.Config{}
	cfg.Repositories = make(map[string][]*gc.Repository, 0)
	cfg.Repositories["discovered"] = repos
	cfg.PrintConfig()
	return nil
}
