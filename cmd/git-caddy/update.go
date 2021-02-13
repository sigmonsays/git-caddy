package main

import (
	gc "github.com/sigmonsays/git-caddy"
)

func UpdateRepo(cfg *gc.Config, repo *gc.Repository, done func()) error {
	log.Debugf("Updating repo %s, remote:%s ", repo.Name, repo.Remote)
	defer done()

	pull := &Pull{cfg, repo}
	err := pull.Run()
	if err != nil {
		return err
	}

	return nil
}
