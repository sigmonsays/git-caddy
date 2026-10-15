package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	gc "github.com/sigmonsays/git-caddy"
)

type Status struct {
	Cfg     *gc.Config
	Repo    *gc.Repository
	Version string
}

func (me *Status) Run(ctx *gc.Context) error {
	cmdline := []string{
		"git",
		"status",
	}
	if me.Cfg.Verbose == false {
		cmdline = append(cmdline, "-s")
	}
	log.Tracef("git status %s", me.Repo.Name)
	log.Tracef("git status command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("status"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("status"))
	c.Dir = ctx.ResolvedWorkingDir
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)
	err := c.Run()
	if err != nil {
		return NewRepoError("Status", me.Repo.Name).WithError(err)
	}

	// git version
	cmdline = []string{
		"git", "describe", "--tags",
	}
	var buf bytes.Buffer
	c = exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = &buf
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("describe"))
	c.Dir = ctx.ResolvedWorkingDir
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)
	// fmt.Printf("Dir %s\n", c.Dir)
	err = c.Run()
	if err == nil {
		me.Version = strings.TrimRight(buf.String(), "\n")
	}

	return nil
}
