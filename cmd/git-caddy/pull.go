package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gc "github.com/sigmonsays/git-caddy"
)

type Pull struct {
	Cfg  *gc.Config
	Repo *gc.Repository
}

func (me *Pull) Run(ctx *gc.Context) error {

	// TODO collect hashes
	branch, err := GetUpstreamBranch(me.Repo.Destination)
	if err == nil {
		ctx.UpstreamBranchName = branch
		log.Tracef("%s upstream branch %s", me.Repo.Name, branch)
	}
	lhash, err := GetLocalHash(me.Repo.Destination, branch)
	if err == nil {
		ctx.LocalHash = lhash
		log.Tracef("%s local hash %s", me.Repo.Name, lhash)
	} else {
		log.Warnf("%s get local hash: %s", me.Repo.Name, err)
	}
	rhash, err := GetRemoteHash(me.Repo.Destination, branch)
	if err == nil {
		ctx.RemoteHash = rhash
		log.Tracef("%s remote hash %s", me.Repo.Name, rhash)
	} else {
		log.Warnf("%s get remote hash: %s", me.Repo.Name, err)
	}

	// perform git pull
	cmdline := []string{
		"git",
		"pull",
		"-q", "--no-edit", "--all",
	}
	log.Tracef("git pull %s", me.Repo.Name)
	err = RepoCommand("pull", me.Cfg, me.Repo, cmdline)
	if err != nil {
		return NewRepoError("Pull", me.Repo.Name).WithError(err)
	}
	return nil
}

func RepoCommand(prefix string, cfg *gc.Config, repo *gc.Repository, cmdline []string) error {
	log.Tracef("repo command: %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, repo.Prefix(prefix))
	c.Stderr = NewPrefixWriter(os.Stderr, repo.Prefix(prefix))
	c.Dir = repo.Destination
	c.Env = populateEnv(c.Env, cfg, repo)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

// get remote branch name
func GetUpstreamBranch(dir string) (string, error) {
	cmdline := []string{
		"ls-remote",
		"-h",
		".",
	}
	c := exec.Command("git", cmdline...)
	c.Dir = dir

	out, err := c.Output()
	if err != nil {
		log.Tracef("get_upstream_branch: [cmdline %s] error: %s\n", cmdline, err)
		return "", nil
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("empty output")
	}
	line := lines[0]
	branch := filepath.Base(line)

	return branch, nil
}

// get a commit hash for branch in git repo at directory
func GetLocalHash(dir, branch string) (string, error) {
	cmdline := []string{
		"-C", dir,
		"rev-list",
		"--max-count=1",
		branch,
	}
	out, err := exec.Command("git", cmdline...).Output()
	if err != nil {
		log.Tracef("local_hash: [cmdline %s] error: %s\n", cmdline, err)
		return "", err
	}
	return strings.Trim(string(out), "\n"), nil
}

func GetRemoteHash(dir, branch string) (string, error) {
	cmdline := []string{
		"-C", dir,
		"ls-remote",
		"origin",
		"-h",
		fmt.Sprintf("refs/heads/%s", branch),
	}
	out, err := exec.Command("git", cmdline...).Output()
	if err != nil {
		log.Tracef("remote_hash: [cmdline %s] error: %s\n", cmdline, err)
		return "", err
	}

	tmp := strings.Fields(string(out))
	if len(tmp) < 1 {
		return "", fmt.Errorf("Invalid input")
	}
	return tmp[0], nil
}
