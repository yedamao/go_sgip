package main

import (
	"flag"
	"log"

	"github.com/yedamao/go_sgip/sgip"
	"github.com/yedamao/go_sgip/sgip/sgiptest"
)

var (
	addr  = flag.String("addr", ":8001", "上行监听地址")
	count = flag.Int("count", 5, "worker 数量")
)

func init() {
	flag.Parse()
}

func main() {

	handler := &sgiptest.MockHandler{}
	receiver, err := sgip.NewReceiver(*addr, *count, handler, false)
	if err != nil {
		log.Println("New Receiver error: ", err)
	}

	HandleSignals(receiver.Stop)

	receiver.Run()

	log.Println("Done")
}
