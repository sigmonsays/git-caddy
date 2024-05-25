package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	for _, name := range names {
		if name.IsDir() == false {
			continue
		}
		repodir := filepath.Join(wd, name.Name())
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
	}
	return nil
}
