package watch

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("watch", func(newlog gologging.Logger) { log = newlog })
}
