package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
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

	nresults := len(me.Repositories)

	me.summary.Results = make([]*ProcessResult, nresults)

	for i, crepo := range me.Repositories {
		wg.Add(1)
		ticket <- true
		n := i + 1
		go ProcessRepo(n, me.Opts, crepo.Cfg, crepo, donefn, me.summary)
	}

	wg.Wait()

	// output table
	switch {
	case me.Opts.StatusOnly && me.Opts.Json:
		err := renderStatusJson(me.summary.Results)
		if err != nil {
			return err
		}
	case me.Opts.StatusOnly:
		renderStatusTable(me.summary.Results)
	case me.Opts.Json:
		err := renderRunJson(me.summary.Results)
		if err != nil {
			return err
		}
	default:
		renderRunTable(me.summary.Results)
	}

	// output any errors
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

func renderRunTable(results []*ProcessResult) {
	tb := tablewriter.NewWriter(os.Stdout)
	header := []string{
		"REPO",
		"VERSION",
	}
	tb.Header(header)
	for _, res := range results {
		row := []string{
			res.Repo.Repo.Name,
			res.Version,
		}
		tb.Append(row)
	}
	tb.Render()
}

// a single row of run output
type RunRow struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Changed bool   `json:"changed"`
	Error   string `json:"error,omitempty"`
}

func renderRunJson(results []*ProcessResult) error {
	rows := make([]*RunRow, 0, len(results))
	for _, res := range results {
		row := &RunRow{
			Name:    res.Repo.Repo.Name,
			Version: res.Version,
			Changed: res.Changed,
		}
		if res.Err != nil {
			row.Error = res.Err.Error()
		}
		rows = append(rows, row)
	}
	return writeJson(rows)
}

func writeJson(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// a single row of status output
type StatusRow struct {
	Name    string `json:"name"`
	Branch  string `json:"branch"`
	Version string `json:"version"`
	State   string `json:"state"`
	Changes int    `json:"changes"`

	// nil if there is no upstream branch
	Ahead  *int `json:"ahead"`
	Behind *int `json:"behind"`

	Error string `json:"error,omitempty"`
}

func NewStatusRow(res *ProcessResult) *StatusRow {
	row := &StatusRow{Name: res.Repo.Repo.Name}
	st := res.Status
	switch {
	case res.Err != nil:
		row.State = "error"
		row.Error = res.Err.Error()
	case res.Exists == false:
		row.State = "missing"
	case st != nil:
		row.Branch = st.Branch
		row.Version = st.Version
		row.State = st.State()
		row.Changes = st.Changes
		if st.HasUpstream {
			row.Ahead = &st.Ahead
			row.Behind = &st.Behind
		}
	}
	return row
}

func renderStatusJson(results []*ProcessResult) error {
	rows := make([]*StatusRow, 0, len(results))
	for _, res := range results {
		rows = append(rows, NewStatusRow(res))
	}
	return writeJson(rows)
}

func renderStatusTable(results []*ProcessResult) {
	tb := tablewriter.NewWriter(os.Stdout)
	header := []string{
		"REPO",
		"BRANCH",
		"VERSION",
		"STATE",
		"CHANGES",
		"AHEAD",
		"BEHIND",
	}
	tb.Header(header)
	for _, res := range results {
		r := NewStatusRow(res)
		row := []string{r.Name, r.Branch, r.Version, r.State, "", "", ""}
		if r.State == "clean" || r.State == "modified" {
			row[4] = strconv.Itoa(r.Changes)
		}
		if r.Ahead != nil {
			row[5] = strconv.Itoa(*r.Ahead)
			row[6] = strconv.Itoa(*r.Behind)
		}
		tb.Append(row)
	}
	tb.Render()
}

type ProcessResult struct {
	Repo    *CompiledRepository
	Version string

	// true if a pull changed the repository
	Changed bool

	// populated in status only mode
	Exists bool
	Status *RepoStatus
	Err    error
}

func ProcessRepo(jobid int, opts *Options, cfg *gc.Config, crepo *CompiledRepository, done func(error), summary *RunSummary) (err error) {
	summary.IncrScanned()
	repo := crepo.Repo

	res := &ProcessResult{
		Repo: crepo,
	}
	defer summary.SaveResult(jobid, res)

	// create new log object
	lvl, _ := gologging.LevelFromString(log.GetLevel())
	dlog := gologging.NewStandardLogger3(lvl, 5)
	rlog := gologging.NewPrefixLogger(fmt.Sprintf("job%d %s: ", jobid, crepo.Repo.Name), dlog)

	rlog.Debugf("Updating repo %s, remote:%s (workingdir:%s) ",
		repo.Name, repo.Remote, crepo.WorkingDir)
	defer func() {
		res.Err = err
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
	res.Exists = repoExists

	if opts.StatusOnly {
		if repoExists == false {
			return nil
		}
		res.Status, err = CollectRepoStatus(cfg, repo, destination)
		if err != nil {
			return err
		}
		res.Version = res.Status.Version
		return nil
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
			res.Changed = pres.Changed
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

	status := &Status{
		Cfg:  cfg,
		Repo: repo,
	}
	err = status.Run(ctx)
	if err != nil {
		return err
	}

	res.Version = status.Version

	rlog.Tracef("UpdateRepo %s: finished without error", repo.Name)
	return nil
}
