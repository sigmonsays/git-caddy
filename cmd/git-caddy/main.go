package main

import (
	"flag"
	"os"
	"sync"
	"time"

	gc "github.com/sigmonsays/git-caddy"
	gologging "github.com/sigmonsays/go-logging"
)

func main() {
	loglevel := "info"
	section := ""
	configfile := "repositories.yaml"
	workingdir := ""
	updateInterval := 0
	flag.StringVar(&configfile, "c", configfile, "specify config file")
	flag.StringVar(&section, "s", section, "section in config file")
	flag.StringVar(&workingdir, "W", workingdir, "initial working directory")
	flag.StringVar(&loglevel, "loglevel", loglevel, "log level")
	flag.StringVar(&loglevel, "l", loglevel, "short for -loglevel")
	flag.IntVar(&updateInterval, "I", updateInterval, "pull upstream for changes")
	flag.Parse()

	gologging.SetLogLevel(loglevel)

	cfg := &gc.Config{}

	if workingdir != "" {
		err := os.Chdir(workingdir)
		ExitIfError(err, "Chdir %s: %s", workingdir, err)
	}

	err := cfg.LoadYaml(configfile)
	ExitIfError(err, "LoadYaml %s: %s", configfile, err)

	if log.IsTrace() {
		cfg.PrintConfig()
	}

	repos, found := cfg.Repositories[section]
	if found == false {
		ExitError("Section not found: %q", section)
	}
	log.Debugf("concurrency:%d", cfg.Concurrency)

	updateRun := &UpdateRepositories{
		Section:      section,
		Cfg:          cfg,
		Repositories: repos,
	}

	if updateInterval == 0 {
		err = updateRun.Run()
		ExitIfError(err, "%s", err)
	}

	tick := time.NewTicker(time.Duration(updateInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err = updateRun.Run()
			if err != nil {
				log.Warnf("%s", err)
			}
		}
	}

}

type UpdateRepositories struct {
	Section      string
	Cfg          *gc.Config
	Repositories []*gc.Repository
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
		n = i + 1
		err := repo.Validate()
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
		go UpdateRepo(me.Cfg, repo, donefn)
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
