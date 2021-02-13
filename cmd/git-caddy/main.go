package main

import (
	"flag"
	"sync"

	gc "github.com/sigmonsays/git-caddy"
	gologging "github.com/sigmonsays/go-logging"
)

func main() {
	loglevel := "trace"
	section := ""
	configfile := "repositories.yaml"
	flag.StringVar(&configfile, "f", configfile, "specify config file")
	flag.StringVar(&section, "s", section, "section in config file")
	flag.StringVar(&loglevel, "loglevel", loglevel, "log level")
	flag.Parse()

	gologging.SetLogLevel(loglevel)

	cfg := &gc.Config{}

	err := cfg.LoadYaml(configfile)
	ExitIfError(err, "LoadYaml %s: %s", configfile, err)

	if log.IsTrace() {
		cfg.PrintConfig()
	}

	repos, found := cfg.Repositories[section]
	if found == false {
		ExitError("Section not found: %q", section)
	}

	var wg sync.WaitGroup
	donefn := func() {
		wg.Done()
	}

	for _, repo := range repos {
		wg.Add(1)
		go UpdateRepo(cfg, repo, donefn)
	}

	wg.Wait()

}
