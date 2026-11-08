package main

import (
	"io"
	"os"
	"os/exec"

	gc "github.com/sigmonsays/git-caddy"
)

// where stdout of git subprocesses is written. Redirected to stderr when
// outputting json so stdout only contains the json document
var cmdStdout io.Writer = os.Stdout

func RepoCommand(prefix string, cfg *gc.Config, repo *gc.Repository, dir string, cmdline []string) error {
	log.Tracef("repo command: %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(cmdStdout, repo.Prefix(prefix))
	c.Stderr = NewPrefixWriter(os.Stderr, repo.Prefix(prefix))
	c.Dir = dir
	c.Env = populateEnv(c.Env, cfg, repo)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
