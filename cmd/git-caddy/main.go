package main

import (
	"flag"
	"fmt"

	gologging "github.com/sigmonsays/go-logging"
)

func main() {

	loglevel := "info"
	configfile := "/etc/whatever.yaml"
	flagvar := 0
	flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	flag.StringVar(&configfile, "config", configfile, "specify config file")
	flag.StringVar(&loglevel, "loglevel", loglevel, "log level")
	flag.Parse()

	gologging.SetLogLevel(loglevel)

	fmt.Printf("Load config from %s\n", configfile)
}
