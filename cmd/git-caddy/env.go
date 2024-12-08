package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shurcooL/go/osutil"
	gc "github.com/sigmonsays/git-caddy"
)

func populateEnv(e []string, cfg *gc.Config, r *gc.Repository) []string {

	// set users home
	// homedir := os.Getenv("HOME")
	// e = append(e, fmt.Sprintf("HOME=%s", homedir))

	gitconfigFile := "/home/sig/.git-caddy.gitconfig"
	e = append(e, fmt.Sprintf("GIT_CONFIG_GLOBAL=%s", gitconfigFile))

	if r.IdentityFile != "" {
		e = env_ssh_command(cfg, r, e, r.IdentityFile)
		return e
	}

	// Check if there is a matching identities block
	if len(cfg.Identities) == 0 {
		return e
	}
	log.Tracef("Looking through %d identities for a match", len(cfg.Identities))
	idmap := make(map[string]*gc.Identity, 0)
	for _, ident := range cfg.Identities {
		for _, repo := range ident.Repositories {
			idmap[repo] = ident
		}
	}
	identity, found := idmap[r.Section]
	if found == false {
		return e
	}
	log.Tracef("Found identity configure for section %s, identity_file %s",
		r.Section, identity.IdentityFile)
	e = env_ssh_command(cfg, r, e, identity.IdentityFile)

	// populate the name and e-mail (for commit messages)
	e = env_author(cfg, r, e, identity)

	log.Tracef("env %+v", e)
	return e
}

func env_author(
	cfg *gc.Config,
	r *gc.Repository,
	e []string,
	id *gc.Identity,
) []string {

	econf := map[string]string{
		"user.name":  id.FullName,
		"user.email": id.Email,
	}

	GitCfgCounter := 0 // todo: Need to pass this in or track it somewhere
	idx := GitCfgCounter
	env := osutil.Environ(e)
	for k, v := range econf {
		keyenv := fmt.Sprintf("GIT_CONFIG_KEY_%d", idx)
		valenv := fmt.Sprintf("GIT_CONFIG_VALUE_%d", idx)
		env.Set(keyenv, k)
		env.Set(valenv, v)
		idx += 1
	}
	cntenv := fmt.Sprintf("%d", idx)
	env.Set("GIT_CONFIG_COUNT", cntenv)
	return env
}
func env_ssh_command(
	cfg *gc.Config,
	r *gc.Repository,
	e []string,
	identityFile string,
) []string {

	sshbin, err := exec.LookPath("ssh")
	if err != nil {
		sshbin = "ssh"
	}
	ssh_opts := makeControlSocket(cfg, r, identityFile)
	if identityFile != "" {
		ssh_opts += " -i " + identityFile
	}
	ssh_command := sshbin + " " + strings.Trim(ssh_opts, " ")
	log.Tracef("repo %s: setting GIT_SSH_COMMAND to %q", r.Name, ssh_command)
	env := osutil.Environ(e)
	env.Set("GIT_SSH_COMMAND", ssh_command)
	return env
}

func makeControlSocket(cfg *gc.Config, r *gc.Repository, ident string) string {
	// todo: Do a better job with the identity file
	b := filepath.Base(ident)
	// maybe IdentitiesOnly=yes  ?
	ret := " -oControlMaster=auto "
	ret += " -oControlPersist=yes "
	ret += " -oControlPath=/tmp/ssh-git-caddy-%u-%h-%n-%p-" + b
	return ret
}
