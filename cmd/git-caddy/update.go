package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	gc "github.com/sigmonsays/git-caddy"
)

type UpdateRepositories struct {
	Section      string
	Cfg          *gc.Config
	Repositories []*gc.Repository

	summary *RunSummary
}

func (me *UpdateRepositories) Run() error {
	var doneMx sync.Mutex
	var errors []error
	ticket := make(chan bool, me.Cfg.Concurrency)
	var wg sync.WaitGroup
	donefn := func(err error) {
		<-ticket
		if err != nil {
			doneMx.Lock()
			errors = append(errors, err)
			doneMx.Unlock()
		}
		wg.Done()
	}

	var n int
	for i, repo := range me.Repositories {

		if repo.Section == "" {
			repo.Section = me.Section
		}
		err := repo.Defaults()
		if err != nil {
			log.Debugf("repo #%d: %s failed setting defaults: %s", n, repo.Name, err)
		}
		n = i + 1
		err = repo.Validate()
		if err != nil {
			log.Warnf("repo #%d: %s failed validation: %s", n, repo.Name, err)
			continue
		}
		if repo.IsEnabled() == false {
			log.Debugf("repo %s is disabled", repo.Name)
			continue
		}
		wg.Add(1)
		ticket <- true
		go UpdateRepo(me.Cfg, repo, donefn, me.summary)
	}

	wg.Wait()

	if len(errors) == 0 {
		log.Debugf("Finished with no errors")
	} else {
		log.Warnf("%d errors occurred", len(errors))
		for i, err := range errors {
			log.Warnf("error #%d: %s", i+1, err)
		}
	}
	return nil
}

func UpdateRepo(cfg *gc.Config, repo *gc.Repository, done func(error), summary *RunSummary) (err error) {
	summary.IncrScanned()

	log.Debugf("Updating repo %s, remote:%s ", repo.Name, repo.Remote)
	defer func() {
		done(err)
		if err != nil {
			summary.IncrErrors()
		}
	}()
	repoExists := false
	isDir := false
	st, err := os.Stat(repo.Destination)
	if err == nil {
		repoExists = true
		isDir = st.IsDir()
	}
	log.Tracef("stat %s; isdir:%v", repo.Destination, isDir)
	if err == nil && isDir == false {
		return fmt.Errorf("%s is not a directory", repo.Destination)
	}

	log.Tracef("repo:%s destination:%s repoExists:%v noClone:%v",
		repo.Name, repo.Destination, repoExists, repo.NoClone)
	if repoExists == false && repo.NoClone == false {
		clone := &Clone{cfg, repo}
		err = clone.Run()
		if err != nil {
			return err
		}
	}

	if strings.Trim(repo.AddFiles, " ") != "" {
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

	status := &Status{cfg, repo}
	err = status.Run()
	if err != nil {
		return err
	}

	log.Tracef("UpdateRepo %s: finished without error", repo.Name)
	return nil
}
