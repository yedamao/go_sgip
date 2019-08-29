package sgiptest

import (
	"bytes"
	"log"
	"net"

	"github.com/yedamao/encoding"
	connp "github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

func newServerSession(rawConn net.Conn) {
	s := &serverSession{*connp.NewConn(rawConn)}
	go s.start()
}

// 代表sp->运营商的一条连接
type serverSession struct {
	connp.Conn
}

func (s *serverSession) BindResp(seq [3]uint32, status protocol.RespStatus) error {
	op, err := protocol.NewResponse(protocol.SGIP_BIND_REP, seq, status)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *serverSession) UnBindResp(seq [3]uint32, status protocol.RespStatus) error {
	op, err := protocol.NewResponse(protocol.SGIP_UNBIND_REP, seq, status)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *serverSession) SubmitResp(seq [3]uint32, status protocol.RespStatus) error {
	op, err := protocol.NewResponse(protocol.SGIP_SUBMIT_REP, seq, status)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *serverSession) start() {
	defer s.Close()

	for {
		op, err := s.Read()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Println(err)
			return
		}

		log.Println(op)

		switch op.GetHeader().CmdId {
		case protocol.SGIP_BIND:
			// check bind
			stat := s.handleBind(op.(*protocol.Bind))
			s.BindResp(op.GetHeader().Sequence, protocol.STAT_OK)
			if stat != protocol.STAT_OK {
				return
			}

		case protocol.SGIP_SUBMIT:
			// check submit
			s.SubmitResp(op.GetHeader().Sequence, protocol.STAT_OK)

		case protocol.SGIP_UNBIND:
			s.UnBindResp(op.GetHeader().Sequence, protocol.STAT_OK)
			return

		default:
			log.Println("not support CmdId. close session.")
			return
		}
	}
}

func (s *serverSession) handleBind(op *protocol.Bind) protocol.RespStatus {
	log.Println("---- handle bind ----")
	log.Println("Type: ", op.Type)
	log.Println("Name: ", op.Name)
	log.Println("Password: ", op.Password)

	if op.Name.String() != "fakename" || op.Password.String() != "1234" {
		log.Println("name/password not match")
		return protocol.STAT_ILLLOGIN
	}

	if op.Type != 1 { // 登录类型 1 sp -> SMG, 2 SMG -> SP
		log.Println("login type is wrong")
		return protocol.STAT_ERLGNTYPE
	}

	return protocol.STAT_OK
}

func (s *serverSession) handleSubmit(op *protocol.Submit) {

	log.Println("SPNumber: ", op.SPNumber)
	log.Println("UserNumber: ", op.UserNumber)
	log.Println("TP_pid: ", op.TP_pid)
	log.Println("TP_udhi: ", op.TP_udhi)
	log.Println("MessageCoding: ", op.MessageCoding)

	messageContent := op.MessageContent.Byte()
	var msg []byte
	switch op.MessageCoding {
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
}
