package main

import (
	"os"
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

type AddFiles struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *AddFiles) Run() error {
	cmdline := []string{
		"git",
		"add",
	}
	log.Tracef("git add %s: %s", me.Repo.Name, me.Repo.AddFiles)

	cmdline = append(cmdline, me.Repo.AddFiles)
	log.Tracef("git add command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = me.Repo.Destination
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
