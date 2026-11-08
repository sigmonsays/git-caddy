package main

import (
	"os/exec"
	"strconv"
	"strings"

	gc "github.com/sigmonsays/git-caddy"
)

// status of a repository, collected without modifying it
type RepoStatus struct {
	Branch  string
	Version string
	Changes int // number of changed/untracked files
	Ahead   int
	Behind  int

	// true if the branch has an upstream configured
	HasUpstream bool
}

func (me *RepoStatus) State() string {
	if me.Changes > 0 {
		return "modified"
	}
	return "clean"
}

// collect status for repository at dir. Only local state is examined; no fetch is performed
func CollectRepoStatus(cfg *gc.Config, repo *gc.Repository, dir string) (*RepoStatus, error) {
	ret := &RepoStatus{}

	out, err := repoOutput(cfg, repo, dir, "status", "--porcelain")
	if err != nil {
		return nil, NewRepoError("Status", repo.Name).WithError(err)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) != "" {
			ret.Changes++
		}
	}

	out, err = repoOutput(cfg, repo, dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err == nil {
		ret.Branch = strings.TrimSpace(out)
	}

	out, err = repoOutput(cfg, repo, dir, "describe", "--tags")
	if err == nil {
		ret.Version = strings.TrimSpace(out)
	}

	out, err = repoOutput(cfg, repo, dir, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if err == nil {
		tmp := strings.Fields(out)
		if len(tmp) == 2 {
			ret.HasUpstream = true
			ret.Ahead, _ = strconv.Atoi(tmp[0])
			ret.Behind, _ = strconv.Atoi(tmp[1])
		}
	}

	return ret, nil
}

func repoOutput(cfg *gc.Config, repo *gc.Repository, dir string, args ...string) (string, error) {
	cmdline := append([]string{"-C", dir}, args...)
	log.Tracef("repo output: git %v", cmdline)
	c := exec.Command("git", cmdline...)
	c.Dir = dir
	c.Env = populateEnv(c.Env, cfg, repo)
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
