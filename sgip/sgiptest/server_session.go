package sgiptest

import (
	"log"
	"net"

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
			// TODO check bind
			s.BindResp(op.GetHeader().Sequence, protocol.STAT_OK)

		case protocol.SGIP_SUBMIT:
			// TODO check submit
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
