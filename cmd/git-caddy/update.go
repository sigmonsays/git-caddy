package main

import (
	"os"

	gc "github.com/sigmonsays/git-caddy"
)

func UpdateRepo(cfg *gc.Config, repo *gc.Repository, done func()) error {
	log.Debugf("Updating repo %s, remote:%s ", repo.Name, repo.Remote)
	defer done()

	repoExists := false
	err := os.Chdir(repo.Destination)
	if err == nil {
		repoExists = true
	}

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
	}

	return nil
}
