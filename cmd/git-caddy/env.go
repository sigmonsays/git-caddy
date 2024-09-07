package main

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shurcooL/go/osutil"
	gc "github.com/sigmonsays/git-caddy"
)

func env_ssh_command(
	cfg *gc.Config,
	r *gc.Repository,
	e []string,
	identityFile string) []string {

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
	env := osutil.Environ([]string{})
	env.Set("GIT_SSH_COMMAND", ssh_command)
	return env
}

func makeControlSocket(cfg *gc.Config, r *gc.Repository, ident string) string {
	// todo: Do a better job with the identity file
	b := filepath.Base(ident)
	ret := " -oControlMaster=auto "
	ret += " -oControlPersist=yes "
	ret += " -oControlPath=/tmp/ssh-git-caddy-%u-%h-%n:%p-" + b
	return ret
}
func populateEnv(e []string, cfg *gc.Config, r *gc.Repository) []string {

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
	log.Tracef("env %+v", e)
	return e
}
