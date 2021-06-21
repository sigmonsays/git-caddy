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

func (me *Commit) ChangedLocally() (bool, error) {

	cmdline := []string{
		"git", "diff-index", "--quiet", "HEAD", "--",
	}
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("commit"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("commit"))
	c.Dir = me.Repo.Destination
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)
	err := c.Run()
	if err != nil {
		var exitcode int
		if c.ProcessState != nil {
			exitcode = c.ProcessState.ExitCode()
		}
		if exitcode == 1 {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (me *Commit) Run() error {

	changed, err := me.ChangedLocally()
	if err != nil {
		return err
	}

	if changed == false {
		log.Tracef("nothing changed in repo %s", me.Repo.Name)
		return nil
	}

	cmdline := []string{
		"git",
		"commit",
		"-m",
		"git-caddy auto commit",
	}
	if me.Cfg.Verbose == false {
		cmdline = append(cmdline, "-q")
	}
	log.Tracef("git commit %s", me.Repo.Name)
	log.Tracef("git commit command %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, me.Repo.Prefix("commit"))
	c.Stderr = NewPrefixWriter(os.Stderr, me.Repo.Prefix("commit"))
	c.Dir = me.Repo.Destination
	c.Env = populateEnv(c.Env, me.Cfg, me.Repo)
	err = c.Run()
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
