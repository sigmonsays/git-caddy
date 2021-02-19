package main

import (
	"fmt"
	"os"
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

type Clone struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Clone) Run() error {
	cmdline := []string{
		"git",
		"--no-pager",
		"clone",
	}
	log.Tracef("git clone repo %s, %s to %s",
		me.Repo.Name, me.Repo.Remote, me.Repo.Destination)

	cmdline = append(cmdline, me.Repo.Remote)
	cmdline = append(cmdline, me.Repo.Destination)
	if me.Repo.Depth > 0 {
		cmdline = append(cmdline, "--depth",
			fmt.Sprintf("%d", me.Repo.Depth))
	}

	log.Tracef("git clone command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("clone"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("clone"))
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)

	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
