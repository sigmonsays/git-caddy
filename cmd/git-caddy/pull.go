package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	gc "github.com/sigmonsays/git-caddy"
)

type Pull struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Pull) Run(ctx *gc.Context) error {

	// TODO collect hashes
	// lhash, err := GetLocalHash(dir string, branch string)

	// perform git pull
	cmdline := []string{
		"git",
		"pull",
		"-q", "--no-edit", "--all",
	}
	log.Tracef("git pull %s", me.Repo.Name)

	log.Tracef("git pull command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("pull"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("pull"))
	c.Dir = me.Repo.Destination
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)
	err := c.Run()
	if err != nil {
		return NewRepoError("Pull", me.Repo.Name).WithError(err)
	}
	return nil
}

// get a commit hash for branch in git repo at directory
func GetLocalHash(dir, branch string) string {
	cmdline := []string{
		"-C", dir,
		"rev-list",
		"--max-count=1",
		branch,
	}
	out, err := exec.Command("git", cmdline...).Output()
	if err != nil {
		log.Infof("remote_hash: [cmdline %s] error: %s\n", cmdline, err)
		return ""
	}
	return strings.Trim(string(out), "\n")
}

func GetRemoteHash(dir, branch string) string {
	cmdline := []string{
		"-C", dir,
		"ls-remote",
		"origin",
		"-h",
		fmt.Sprintf("refs/heads/%s", branch),
	}
	out, err := exec.Command("git", cmdline...).Output()
	if err != nil {
		log.Infof("remote_hash: [cmdline %s] error: %s\n", cmdline, err)
		return ""
	}
	tmp := strings.Fields(string(out))
	if len(tmp) < 1 {
		return ""
	}
	return tmp[0]
}
