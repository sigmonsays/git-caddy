package main

import (
	"os"
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
		"-q", "--no-edit", "--all",
	}
	log.Tracef("git pull %s", me.Repo.Name)

	cmdline = append(cmdline, me.Repo.Remote)
	log.Tracef("git pull command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("pull"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("pull"))
	c.Dir = me.Repo.Destination
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
