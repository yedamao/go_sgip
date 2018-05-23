package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func HandleSignals(stopFunction func()) {
	var callback sync.Once

	// On ^C or SIGTERM, gracefully stop the sniffer
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		log.Println("service", "Received sigterm/sigint, stopping")
		callback.Do(stopFunction)
	}()
}
