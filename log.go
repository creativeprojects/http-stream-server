package main

import "github.com/creativeprojects/clog"

// setupLogger returns a cleaning function, if any
func setupLogger(flags Flags) func() {
	level := clog.LevelInfo
	if flags.debug {
		level = clog.LevelTrace
	} else if flags.verbose {
		level = clog.LevelDebug
	} else if flags.quiet {
		level = clog.LevelWarning
	}
	clog.SetDefaultLogger(clog.NewFilteredConsoleLogger(level))
	return nil
}
