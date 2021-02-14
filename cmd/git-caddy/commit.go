package main

import (
	"os"
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

type Commit struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Commit) Run() error {
	cmdline := []string{
		"git",
		"commit",
		"-m",
		"git-caddy auto commit",
	}
	log.Tracef("git commit %s", me.Repo.Name)
	log.Tracef("git commit command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix())
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix())
	c.Dir = me.Repo.Destination
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
