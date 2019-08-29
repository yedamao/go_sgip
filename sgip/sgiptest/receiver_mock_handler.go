package sgiptest

import (
	"bytes"
	"log"

	"github.com/yedamao/encoding"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

type MockHandler struct{}

func (h *MockHandler) OnBind(login_type uint8, name, password string) protocol.RespStatus {
	log.Println("---- handle bind ----")
	log.Println("Type: ", login_type)
	log.Println("Name: ", name)
	log.Println("Password: ", password)

	if name != "fakename" || password != "1234" {
		log.Println("name/password not match")
		return protocol.STAT_ILLLOGIN
	}

	if login_type != 2 { // 登录类型 1 sp -> SMG, 2 SMG -> SP
		log.Println("login type is wrong")
		return protocol.STAT_ERLGNTYPE
	}

	return protocol.STAT_OK
}

func (h *MockHandler) OnDeliver(
	userNumber, spNumber string, TP_pid, TP_udhi,
	messageCoding uint8, messageContent []byte,
) protocol.RespStatus {
	log.Println("---- handle deliver ---- ")
	log.Println("UserNumber: ", userNumber)
	log.Println("SPNumber: ", spNumber)
	log.Println("TP_pid: ", TP_pid)
	log.Println("TP_udhi: ", TP_udhi)
	log.Println("MessageCoding: ", messageCoding)

	var msg []byte
	switch messageCoding {
	case protocol.ASCII:
		msg = messageContent
	case protocol.GBK:
		// convert GB to UTF-8
		msg = encoding.GB18030(messageContent).Decode()
	case protocol.UCS2:
		// 0x05 0x00 0x03  长短信
		if bytes.Equal(messageContent[:3], []byte{0x05, 0x00, 0x03}) {
			// convert UCS2 to UTF-8 砍掉六字节长短信头
			msg = encoding.UCS2(messageContent[6:]).Decode()
		} else {
			// convert UCS2 to UTF-8
			msg = encoding.UCS2(messageContent).Decode()
		}
	}

	log.Println("MessageContent: ", string(msg))

	return protocol.STAT_OK
}

func (h *MockHandler) OnReport(
	seq [3]uint32, reportType uint8, userNumber string,
	state, errorCode uint8,
) protocol.RespStatus {
	log.Println("---- handle report ---- ")
	log.Println("Sequence: ", seq)
	log.Println("ReportType: ", reportType)
	log.Println("UserNumber: ", userNumber)
	log.Println("State: ", state)
	log.Println("ErrorCode: ", errorCode)

	return protocol.STAT_OK
}
