package main

import (
	"fmt"
	"os"

	gc "github.com/sigmonsays/git-caddy"
)

func UpdateRepo(cfg *gc.Config, repo *gc.Repository, done func(error)) (err error) {
	log.Debugf("Updating repo %s, remote:%s ", repo.Name, repo.Remote)
	defer done(err)
	repoExists := false
	isDir := false
	st, err := os.Stat(repo.Destination)
	if err == nil {
		repoExists = true
		isDir = st.IsDir()
	}
	log.Tracef("stat %s; isdir:%v", repo.Destination, isDir)
	if isDir == false {
		return fmt.Errorf("%s is not a directory", repo.Destination)
	}

	log.Tracef("repo:%s destination:%s repoExists:%v",
		repo.Name, repo.Destination, repoExists)
	if repoExists == false {
		clone := &Clone{cfg, repo}
		err = clone.Run()
		if err != nil {
			return err
		}
	}

	if repo.AddFiles != "" {
		addFiles := &AddFiles{cfg, repo}
		err = addFiles.Run()
		if err != nil {
			return err
		}
	}

	if repoExists == true {
		pull := &Pull{cfg, repo}
		err = pull.Run()
		if err != nil {
			return err
		}

		commit := &Commit{cfg, repo}
		err = commit.Run()
		if err != nil {
			return err
		}

		push := &Push{cfg, repo}
		err = push.Run()
		if err != nil {
			return err
		}
	}

	log.Tracef("UpdateRepo %s: finished without error", repo.Name)
	return nil
}
