package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/magicLian/gostarter/pkg/setting"
)

func main() {
	var (
		configFile = flag.String("config", "", "path to config file")
		homePath   = flag.String("homepath", "", "path to go starter install/home path, defaults to working directory")
	)
	flag.Parse()

	server, err := InitGoStarterWire(&setting.CommandLineArgs{
		Config:   *configFile,
		HomePath: *homePath,
		Args:     flag.Args(),
	})
	if err != nil {
		panic(err.Error())
	}
	go listenToSystemSignals(server)

	if err := server.Run(); err != nil {
		server.Exit(err)
	}
}

func listenToSystemSignals(server *PatrolServer) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:

		case sig := <-signalChan:
			server.Shutdown(fmt.Sprintf("System signal: %s", sig))
		}
	}
}
