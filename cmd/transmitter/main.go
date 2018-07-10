package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/yedamao/encoding"
	"github.com/yedamao/go_sgip/sgip"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

var (
	host        = flag.String("host", "localhost", "SMSC host")
	port        = flag.Int("port", 8801, "SMSC port")
	areaCode    = flag.String("area-code", "010", "长途区号")
	corpId      = flag.String("corp-id", "00000", "5位企业代码")
	serviceType = flag.String("service-type", "", "业务代码，由SP定义")
	name        = flag.String("name", "", "Login Name")
	passwd      = flag.String("passwd", "", "Login Password")

	spNumber   = flag.String("sp-number", "", "SP的接入号码")
	destNumber = flag.String("dest-number", "", "接收手机号码, 86..., 多个使用，分割")
	msg        = flag.String("msg", "", "短信内容")
)

func init() {
	flag.Parse()
}

func main() {

	destNumbers := strings.Split(*destNumber, ",")
	fmt.Println("destNumbers: ", destNumbers)

	tx, err := sgip.NewTransmitter(*host, *port, *areaCode, *corpId, *name, *passwd)
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

	var seq [3]uint32
	// Send sms
	fmt.Println("----- single msg -----")
	seq, err = tx.Submit(*spNumber, destNumbers, *serviceType, 0, 15, content)
	if err != nil {
		fmt.Println("SubmitSm err:", err)
	}

	// Should save this to match with message_id
	fmt.Println("seq:", seq)

	for {
		op, err := tx.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}
		fmt.Println(op)

		switch op.GetHeader().CmdId {
		case protocol.SGIP_SUBMIT_REP:
			// message_id should match this with seq message
			fmt.Println("MSG ID:", op.GetHeader().Sequence)
			tx.Unbind()
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
