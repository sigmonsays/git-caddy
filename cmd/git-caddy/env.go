package main

import (
	"fmt"

	gc "github.com/sigmonsays/git-caddy"
)

func populateEnv(e []string, cfg *gc.Config, r *gc.Repository) []string {

	// see if we need to set the GIT_SSH_COMMAND for a custom identity
	// IdentityFile overides a higher level config
	if r.IdentityFile != "" {
		ssh_command := fmt.Sprintf("ssh -i %s", r.IdentityFile)
		e = append(e, "GIT_SSH_COMMAND="+ssh_command)
		return e
	}

	// Check if there is a matching identities block
	idmap := make(map[string]*gc.Identity, 0)
	for _, ident := range cfg.Identities {
		for _, repo := range ident.Repositories {
			idmap[repo] = ident
		}
	}

	identity, found := idmap[r.section]
	if found == false {
		return e
	}

	return e
}
