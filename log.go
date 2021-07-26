package gitcaddy

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("gitcaddy", func(newlog gologging.Logger) { log = newlog })
}
