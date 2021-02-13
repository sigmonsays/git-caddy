package main

import (
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

type Pull struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Pull) Run() error {
	cmdline := []string{
		"git",
		"pull",
	}
	log.Tracef("git pull %s", me.Repo.Name)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
