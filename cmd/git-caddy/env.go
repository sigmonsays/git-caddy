package main

import (
	"fmt"

	gc "github.com/sigmonsays/git-caddy"
)

func env_GIT_SSH_COMMAND(e []string, identityFile string) []string {
	ssh_command := fmt.Sprintf("ssh -i %s", identityFile)
	e = append(e, "GIT_SSH_COMMAND="+ssh_command)
	return e
}

func populateEnv(e []string, cfg *gc.Config, r *gc.Repository) []string {

	// see if we need to set the GIT_SSH_COMMAND for a custom identity
	// IdentityFile overides a higher level config
	if r.IdentityFile != "" {
		e = env_GIT_SSH_COMMAND(e, r.IdentityFile)
		return e
	}

	// Check if there is a matching identities block
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

	e = env_GIT_SSH_COMMAND(e, identity.IdentityFile)
	log.Tracef("env %+v", e)

	return e
}
