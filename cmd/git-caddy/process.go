package main

import (
	"os"
	"strings"
	"sync"

	gc "github.com/sigmonsays/git-caddy"
)

func RunCompiled(opts *Options, cfg *gc.Config, summary *RunSummary, run *CompiledRun) error {
	p := &ProcessRepositories{
		Opts:         opts,
		Cfg:          cfg,
		summary:      summary,
		Repositories: run.List,
	}
	err := p.Run()
	if err != nil {
		return err
	}
	return nil
}

type ProcessRepositories struct {
	Opts         *Options
	Cfg          *gc.Config
	Repositories []*CompiledRepository
	summary      *RunSummary
}

func (me *ProcessRepositories) Run() error {
	var errors []error
	var doneMx sync.Mutex
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

	for _, crepo := range me.Repositories {
		repo := crepo.Repo
		wg.Add(1)
		ticket <- true
		go ProcessRepo(me.Opts, me.Cfg, repo, donefn, me.summary)
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

func ProcessRepo(opts *Options, cfg *gc.Config, repo *gc.Repository, done func(error), summary *RunSummary) (err error) {
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
		return NewRepoErrorf("Update", repo.Name, "%s is not a directory", repo.Destination)
	}

	ctx := &gc.Context{
		RepoExists: repoExists,
	}

	if opts.Pretend {
		log.Infof("Pretend %s (exists:%v)", repo.Name, repoExists)
		return nil
	}

	log.Tracef("repo:%s destination:%s repoExists:%v noClone:%v",
		repo.Name, repo.Destination, repoExists, repo.NoClone)
	if repoExists == false && repo.NoClone == false {
		clone := &Clone{cfg, repo}
		err = clone.Run(ctx)
		if err != nil {
			return err
		}
	}

	if strings.Trim(repo.AddFiles, " ") != "" {
		addFiles := &AddFiles{cfg, repo}
		err = addFiles.Run(ctx)
		if err != nil {
			return err
		}
	}

	if repoExists == true {
		pull := &Pull{cfg, repo}
		pres, err := pull.Run(ctx)
		if err == nil {
			if pres.Changed {
				summary.Changed += 1
			}
		} else {
			return err
		}

		commit := &Commit{cfg, repo}
		err = commit.Run(ctx)
		if err != nil {
			return err
		}

		push := &Push{cfg, repo}
		err = push.Run(ctx)
		if err != nil {
			return err
		}
	}

	status := &Status{cfg, repo}
	err = status.Run(ctx)
	if err != nil {
		return err
	}

	log.Tracef("UpdateRepo %s: finished without error", repo.Name)
	return nil
}
