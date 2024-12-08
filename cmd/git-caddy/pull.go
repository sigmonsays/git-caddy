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
type PullResult struct {

	// true if the pull changed local files
	Changed bool
}

func (me *Pull) Run(ctx *gc.Context) (*PullResult, error) {

	ret := &PullResult{}

	// collect hashes
	branch, err := GetUpstreamBranch(me.Cfg, me.Repo, ctx.ResolvedWorkingDir)
	if err == nil {
		ctx.UpstreamBranchName = branch
		log.Tracef("%s upstream branch %s", me.Repo.Name, branch)
	}
	lhash, err := GetLocalHash(me.Cfg, me.Repo, ctx.ResolvedWorkingDir, branch)
	if err == nil {
		ctx.LocalHash = lhash
		log.Tracef("%s local hash %s", me.Repo.Name, lhash)
	} else {
		log.Warnf("%s get local hash: %s", me.Repo.Name, err)
	}
	rhash, err := GetRemoteHash(me.Cfg, me.Repo, ctx.ResolvedWorkingDir, branch)
	if err == nil {
		ctx.RemoteHash = rhash
		log.Tracef("%s remote hash %s", me.Repo.Name, rhash)
	} else {
		log.Warnf("%s get remote hash: %s", me.Repo.Name, err)
	}

	// if local hash and remote hash match, we're up to date!!
	nonEmptyHashes := lhash != "" && rhash != ""
	if nonEmptyHashes && lhash == rhash {
		log.Tracef("already up to date")
	} else {
		ret.Changed = true
	}

	// do not pull if the hashes are the same
	if lhash == rhash && lhash != "" {
		return ret, nil
	}

	// perform git pull
	cmdline := []string{
		"git",
		"pull",
		"-q", "--no-edit", "--all",
	}
	log.Tracef("git pull %s", me.Repo.Name)
	err = RepoCommand("pull", me.Cfg, me.Repo, ctx.ResolvedWorkingDir, cmdline)
	if err != nil {
		return nil, NewRepoError("Pull", me.Repo.Name).WithError(err)
	}
	return ret, nil
}

func RepoCommand(prefix string, cfg *gc.Config, repo *gc.Repository, dir string, cmdline []string) error {
	log.Tracef("repo command: %v", cmdline)
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = NewPrefixWriter(os.Stdout, repo.Prefix(prefix))
	c.Stderr = NewPrefixWriter(os.Stderr, repo.Prefix(prefix))
	c.Dir = dir
	c.Env = populateEnv(c.Env, cfg, repo)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

// get remote branch name
func GetUpstreamBranch(cfg *gc.Config, repo *gc.Repository, dir string) (string, error) {
	// git ls-remote https://github.com/git/git.git --symref HEAD
	cmdline := []string{
		"-C",
		dir,
		"ls-remote",
		"--symref",
		repo.Remote,
		"HEAD",
	}
	c := exec.Command("git", cmdline...)
	c.Dir = dir
	c.Env = populateEnv(c.Env, cfg, repo)

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
	tmp := strings.Fields(line)
	if len(tmp) == 0 {
		return "", fmt.Errorf("invalid output")
	}
	branch := filepath.Base(tmp[1])

	return branch, nil
}

// get a commit hash for branch in git repo at directory
func GetLocalHash(cfg *gc.Config, repo *gc.Repository, dir, branch string) (string, error) {
	cmdline := []string{
		"-C", dir,
		"rev-list",
		"--max-count=1",
		branch,
		"--",
	}
	c := exec.Command("git", cmdline...)
	c.Env = populateEnv(c.Env, cfg, repo)
	out, err := c.Output()
	if err != nil {
		log.Tracef("local_hash: [cmdline %s] error: %s\n", cmdline, err)
		return "", err
	}
	return strings.Trim(string(out), "\n"), nil
}

func GetRemoteHash(cfg *gc.Config, repo *gc.Repository, dir, branch string) (string, error) {
	if branch == "" {
		return "", fmt.Errorf("branch empty: branch required")
	}
	cmdline := []string{
		"-C", dir,
		"ls-remote",
		"origin",
		"-h",
		fmt.Sprintf("refs/heads/%s", branch),
	}
	c := exec.Command("git", cmdline...)
	c.Env = populateEnv(c.Env, cfg, repo)
	out, err := c.Output()
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
