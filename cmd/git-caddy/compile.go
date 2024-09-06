package main

import (
	"os"

	gc "github.com/sigmonsays/git-caddy"
)

type Compile struct {
	Section      string
	Cfg          *gc.Config
	Repositories []*gc.Repository
	WorkingDir   string
	summary      *RunSummary
}

type CompiledRepository struct {
	WorkingDir string
	Repo       *gc.Repository
}

func (me *Compile) Run() ([]*CompiledRepository, error) {
	ret := make([]*CompiledRepository, 0)

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

		if len(repo.Names) > 0 {
			// expand names
			for _, name := range repo.Names {
				repo2 := repo.Copy()
				repo2.Name = ""
				repo2.Remote = repo.Remote + name
				repo2.Defaults()
				log.Tracef("expanded repo %s", repo2.Remote)
				crepo2 := &CompiledRepository{
					WorkingDir: os.ExpandEnv(me.WorkingDir),
					Repo:       repo2,
				}
				ret = append(ret, crepo2)
			}
		} else {
			crepo := &CompiledRepository{
				WorkingDir: os.ExpandEnv(me.WorkingDir),
				Repo:       repo,
			}
			ret = append(ret, crepo)
		}
	}
	return ret, nil
}
