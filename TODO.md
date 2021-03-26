# TODO

  
- Implement logging to file on disk
  - optional behavior, log only if parameter is in yaml
  - use rotating logger, ability to set rotate sizes

- show status of repositories and dont update
  git-caddy -X status

- Implement ability to shorten repository specification even more

    repositories_group:
        "git@github.com:sigmonsays/":
            - git-cadd
            - runitcmd
            - screenshot
