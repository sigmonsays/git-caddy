package main

import (
	"os"
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

type Status struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Status) Run() error {
	cmdline := []string{
		"git",
		"status",
	}
	log.Tracef("git status %s", me.Repo.Name)
	log.Tracef("git status command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("status"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("status"))
	c.Dir = me.Repo.Destination
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
