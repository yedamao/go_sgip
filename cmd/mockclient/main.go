package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/yedamao/encoding"
	"github.com/yedamao/go_sgip/sgip/protocol"
	"github.com/yedamao/go_sgip/sgip/sgiptest"
)

var (
	host   = flag.String("host", "localhost", "SP receiver host")
	port   = flag.Int("port", 8001, "SP receiver port")
	name   = flag.String("name", "", "Login Name")
	passwd = flag.String("passwd", "", "Login Password")

	spNumber   = flag.String("sp-number", "", "SP的接入号码")
	userNumber = flag.String("user-number", "", "发送短消息的用户手机号，手机号码前加“86”国别标志")
	msg        = flag.String("msg", "", "短信内容")

	sleep = flag.Int("sleep", 1, "sleep some seconds after receive Deliver response")
)

func init() {
	flag.Parse()
}

func main() {

	client, err := sgiptest.NewSMGClient(*host, *port, *name, *passwd)
	if err != nil {
		fmt.Println("Connection Err:", err)
		return
	}

	fmt.Println("connect succ")

	// encoding msg
	content := encoding.UTF82GBK([]byte(*msg))

	if len(content) > 140 {
		fmt.Println("msg Err: not suport long sms")
	}

	fmt.Println("----- Deliver single msg -----")
	// Send sms
	err = client.Deliver(*userNumber, *spNumber, 0, 0, protocol.GBK, content)
	if err != nil {
		fmt.Println("Deliver: ", err)
	}

	for {
		op, err := client.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}

		fmt.Println(op)

		switch op.GetHeader().CmdId {
		case protocol.SGIP_DELIVER_REP:
			time.Sleep(time.Duration(*sleep) * time.Second)
			client.Unbind()
		case protocol.SGIP_UNBIND_REP:
			fmt.Println("unbind response")
			break
		default:
			fmt.Printf("Unexpect CmdId: %0x\n", op.GetHeader().CmdId)
			fmt.Println("MSG ID:", op.GetHeader().Sequence)
		}
	}

	fmt.Println("ending...")
}
