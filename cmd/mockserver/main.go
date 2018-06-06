package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yedamao/go_sgip/sgip/sgiptest"
)

var (
	addr = flag.String("addr", ":8801", "监听地址")
)

func main() {
	server, err := sgiptest.NewServer(*addr)
	if err != nil {
		flag.Usage()
		os.Exit(-1)
	}

	HandleSignals(server.Stop)

	server.Run()

	fmt.Println("Done")
}

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
