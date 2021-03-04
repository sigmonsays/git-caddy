# TODO

- "No Clone" option
  if no_clone is set, then do not clone the repo if it does not exist
  this makes it useful to checkout a repository manually and have git-caddy update it but nothing else
  
- Implement ability to shorten repository specification even more

    repositories_group:
        "git@github.com:sigmonsays/":
            - git-cadd
            - runitcmd
            - screenshot
