# TODO

- "No Clone" option
  if no_clone is set, then do not clone the repo if it does not exist
  this makes it useful to checkout a repository manually and have git-caddy update it but nothing else
  
- Implement logging to file on disk
  - optional behavior, log only if parameter is in yaml
  - use rotating logger, ability to set rotate sizes

- Allow screenshot name to be given on input
  - prepend/append name to date string?

- show status of repositories and dont update
  git-caddy -X status

- Implement ability to shorten repository specification even more

    repositories_group:
        "git@github.com:sigmonsays/":
            - git-cadd
            - runitcmd
            - screenshot
