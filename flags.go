package main

import (
	"flag"
)

// Flags contains command line flags
type Flags struct {
	configFile string
	quiet      bool
	verbose    bool
	debug      bool
}

var (
	flags Flags
)

func init() {
	flag.StringVar(&flags.configFile, "c", "config.yaml", "configuration file")
	flag.BoolVar(&flags.quiet, "q", false, "quiet - do not send any output")
	flag.BoolVar(&flags.verbose, "v", false, "verbose - display debugging information")
	flag.BoolVar(&flags.debug, "d", false, "debug - display full debugging information")
}
