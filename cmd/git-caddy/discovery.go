package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func execCmd(dir string, cmd []string) (string, error) {
	buf := bytes.NewBuffer(nil)
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = dir
	c.Stdout = buf
	err := c.Run()
	ret := strings.Trim(buf.String(), "\n")
	return ret, err
}
func getRemoteUrl(dir string, origin string) (string, error) {
	getUrlCmd := []string{
		"git", "remote", "get-url", "--push", origin}
	buf, err := execCmd(dir, getUrlCmd)
	if err != nil {
		return "", err
	}
	return buf, nil
}

func getOrigins(dir string) ([]string, error) {
	showOrigin := []string{
		"git", "remote", "show",
	}
	buf, err := execCmd(dir, showOrigin)
	if err != nil {
		return nil, err
	}
	// read first origin
	bio := bytes.NewBufferString(buf)
	ret := make([]string, 0)
	for {
		line, err := bio.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return ret, err
		}
		line = strings.Trim(line, "\n")
		ret = append(ret, line)
	}
	return ret, nil
}

// find first origin with a remote url set
// try 'origin' first to be fast
func findRemoteUrl(dir string) (string, error) {
	// fast path
	remote, err := getRemoteUrl(dir, "origin")
	if err == nil && remote != "" {
		return remote, nil
	}

	// list origins and try each one
	origins, err := getOrigins(dir)
	if err != nil {
		return "", err
	}
	for _, origin := range origins {
		remote, err := getRemoteUrl(dir, origin)
		if err == nil && remote != "" {
			return remote, nil
		}
	}
	return "", fmt.Errorf("Unable to find remote for %s", dir)
}

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
		remote, err := findRemoteUrl(repodir)
		if err != nil {
			log.Debugf("findRemoteUrl %s: %s", repodir, err)
			continue
		}
		if remote == "" {
			log.Debugf("findRemoteUrl %s: remote empty", repodir)
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
