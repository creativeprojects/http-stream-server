package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/creativeprojects/clog"
	"github.com/creativeprojects/http-stream-server/cfg"
	"github.com/creativeprojects/http-stream-server/server"
)

func main() {
	flag.Parse()

	cleanLogger := setupLogger(flags)
	if cleanLogger != nil {
		defer cleanLogger()
	}

	config, err := cfg.LoadFileConfig(flags.configFile)
	if err != nil {
		clog.Errorf("cannot load configuration: %v", err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	servers := setupServers(config)
	if len(servers) > 0 {
		// Wait until we're politely asked to leave
		<-stop
		shutdown(servers)
	}
}

func setupServers(config cfg.Config) map[string]*server.HTTPServer {
	httpServers := make(map[string]*server.HTTPServer, len(config.Servers))
	for name, s := range config.Servers {
		httpServer, err := server.NewHTTPServer(name, s)
		if err != nil {
			clog.Errorf("cannot start server %q: %v", name, err)
			continue
		}
		httpServers[name] = httpServer
		go httpServer.Start()
	}
	return httpServers
}

func shutdown(httpServers map[string]*server.HTTPServer) {
	if len(httpServers) == 0 {
		return
	}
	clog.Info("shutting down...")
	var wg sync.WaitGroup
	wg.Add(len(httpServers))
	for _, s := range httpServers {
		if s == nil {
			wg.Done()
			continue
		}
		go s.Shutdown(&wg, 1*time.Minute)
	}
	wg.Wait()
}
