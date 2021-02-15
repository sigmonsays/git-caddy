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
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("commit"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("commit"))
	c.Dir = me.Repo.Destination
	err := c.Run()
	if err != nil {
		log.Debugf("commit command error : %s", err)

		var exitcode int

		if c.ProcessState != nil {
			exitcode = c.ProcessState.ExitCode()
		}
		if exitcode == 1 {
			// no files to commit
			log.Tracef("treating exit code of 1 for git commit as success")
			return nil
		}
		return err
	}
	return nil
}
