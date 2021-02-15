package main

import (
	"os"
	"os/exec"
	"strings"

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

	cmdline = append(cmdline, strings.Fields(me.Repo.AddFiles)...)
	log.Tracef("git add command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("addfiles"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("addfiles"))
	c.Dir = me.Repo.Destination
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
