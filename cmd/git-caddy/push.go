package main

import (
	"os"
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

type Push struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Push) Run() error {
	cmdline := []string{
		"git",
		"push",
	}
	if me.Cfg.Verbose == false {
		cmdline = append(cmdline, "-q")
	}
	log.Tracef("git push %s", me.Repo.Name)

	log.Tracef("git push command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("push"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("push"))
	c.Dir = me.Repo.Destination
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)
	err := c.Run()
	if err != nil {
		return NewRepoError("Push", me.Repo.Name).WithError(err)
	}
	return nil
}
