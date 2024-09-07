package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	gc "github.com/sigmonsays/git-caddy"
	gologging "github.com/sigmonsays/go-logging"
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

	for i, crepo := range me.Repositories {
		wg.Add(1)
		ticket <- true
		n := i + 1
		go ProcessRepo(n, me.Opts, crepo.Cfg, crepo, donefn, me.summary)
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

func ProcessRepo(jobid int, opts *Options, cfg *gc.Config, crepo *CompiledRepository, done func(error), summary *RunSummary) (err error) {
	summary.IncrScanned()
	repo := crepo.Repo

	// create new log object
	lvl, _ := gologging.LevelFromString(log.GetLevel())
	dlog := gologging.NewStandardLogger3(lvl, 5)
	rlog := gologging.NewPrefixLogger(fmt.Sprintf("job%d %s: ", jobid, crepo.Repo.Name), dlog)

	rlog.Debugf("Updating repo %s, remote:%s ", repo.Name, repo.Remote)
	defer func() {
		done(err)
		if err != nil {
			summary.IncrErrors()
		}
	}()

	// figure out destination directory
	var destination string
	if strings.HasPrefix(repo.Destination, "/") {
		destination = repo.Destination
	} else {
		destination = filepath.Join(crepo.WorkingDir, repo.Destination)
	}
	rlog.Tracef("destination %s", destination)

	repoExists := false
	isDir := false
	st, err := os.Stat(destination)
	if err == nil {
		repoExists = true
		isDir = st.IsDir()
	}
	rlog.Tracef("stat %s; isdir:%v", destination, isDir)
	if err == nil && isDir == false {
		return NewRepoErrorf("Update", repo.Name, "%s is not a directory", destination)
	}

	ctx := &gc.Context{
		RepoExists:         repoExists,
		ResolvedWorkingDir: destination,
	}

	if opts.Pretend {
		rlog.Infof("pretend: repo %s at %s (exists:%v)", repo.Name, destination, repoExists)
		return nil
	}

	rlog.Tracef("repo:%s destination:%s repoExists:%v noClone:%v",
		repo.Name, destination, repoExists, repo.NoClone)
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

	rlog.Tracef("UpdateRepo %s: finished without error", repo.Name)
	return nil
}
