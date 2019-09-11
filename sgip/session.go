package sgip

import (
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

// Session bewteen SP and Operator
type Session struct {
	// SMG -> SP
	conn       conn.Conn
	handler    Handler
	pipe       chan protocol.Operation
	isAuth     bool
	isClose    uint64
	done       chan struct{}
	serverDone chan struct{}

	// debug flag
	debug bool
}

func NewSession(
	connection net.Conn, handler Handler,
	done chan struct{}, debug bool,
) *Session {

	return &Session{
		handler:    handler,
		conn:       *conn.NewConn(connection),
		pipe:       make(chan protocol.Operation),
		done:       make(chan struct{}),
		serverDone: done,
		debug:      debug,
	}
}

func (s *Session) Run() {
	defer s.Close()

	go s.recvWorker()

	for {
		var op protocol.Operation

		select {

		case op = <-s.pipe:
			s.process(op)

		case <-time.After(60 * time.Second):
			// Session connection inactive
			// exceed 60s close session
			log.Println("Session: Conn inactive exceed 60s close session")
			return

		case <-s.done:
			log.Println("Session: To be close")
			return

		case <-s.serverDone:
			log.Println("Session: Receiver server closed")
			return
		}
	}
}

// Close Session
func (s *Session) Close() {
	if atomic.CompareAndSwapUint64(&s.isClose, 0, 1) {
		close(s.done)
		s.conn.Close()
		log.Println("Session: Closing")
		return
	}
}

// recvWorker read Operation on session connection
// pass to session handle loop by pipe channel
func (s *Session) recvWorker() {

	for {

		select {
		case <-s.serverDone:
			log.Println("Session.recvWorker: Receiver server closed")
			return
		case <-s.done:
			log.Println("Session.recvWorker: Session closed")
			return
		default:
		}

		s.conn.SetDeadline(time.Now().Add(1e9))
		op, err := s.conn.Read()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Println("Session.recvWorker: ", err)
			// Close session
			s.Close()
			return
		}

		if s.debug {
			log.Println(op)
		}

		s.pipe <- op
	}
}

// process Operation
func (s *Session) process(op protocol.Operation) {

	switch op.GetHeader().CmdId {
	case protocol.SGIP_BIND:
		bind, ok := op.(*protocol.Bind)
		if !ok {
			log.Println("Session: bind type assert error")
			s.Close()
			return
		}

		// check is authorized
		if s.isAuth {
			s.bindResp(op.GetHeader().Sequence, protocol.STAT_RPTLOGIN)
			break
		}

		stat := s.handler.OnBind(
			bind.Type, bind.Name.String(), bind.Password.String(),
		)

		s.bindResp(op.GetHeader().Sequence, stat)

		// check bin result
		if stat != protocol.STAT_OK {
			log.Println("Session: bind failed")
			s.Close()
			return
		}

	case protocol.SGIP_DELIVER:
		deliver, ok := op.(*protocol.Deliver)
		if !ok {
			log.Println("Session: deliver type assert error")
			s.Close()
			return
		}

		stat := s.handler.OnDeliver(
			deliver.UserNumber.String(), deliver.SPNumber.String(),
			deliver.TP_pid, deliver.TP_udhi, deliver.MessageCoding,
			deliver.MessageContent.Byte(),
		)

		s.deliverResp(op.GetHeader().Sequence, stat)

	case protocol.SGIP_REPORT:
		report, ok := op.(*protocol.Report)
		if !ok {
			log.Println("Session: report type assert error")
			s.Close()
			return
		}

		stat := s.handler.OnReport(
			report.SubmitSequence, report.ReportType, report.UserNumber.String(),
			report.State, report.ErrorCode,
		)

		s.reportResp(op.GetHeader().Sequence, stat)

	case protocol.SGIP_UNBIND:
		s.unbindResp(op.GetHeader().Sequence)
		s.Close()
		return

	default:
		log.Printf("Session: Unknow Operation CmdId: 0x%x\n", op.GetHeader().CmdId)
		s.Close()
		return
	}
}

func (s *Session) bindResp(seq [3]uint32, status protocol.RespStatus) error {
	op, err := protocol.NewResponse(protocol.SGIP_BIND_REP, seq, status)
	if err != nil {
		return err
	}

	return s.conn.Write(op)
}

func (s *Session) unbindResp(seq [3]uint32) error {
	op, err := protocol.NewUnbindResp(seq)
	if err != nil {
		return err
	}

	return s.conn.Write(op)
}

func (s *Session) deliverResp(seq [3]uint32, status protocol.RespStatus) error {
	op, err := protocol.NewResponse(protocol.SGIP_DELIVER_REP, seq, status)
	if err != nil {
		return err
	}

	return s.conn.Write(op)
}

func (s *Session) reportResp(seq [3]uint32, status protocol.RespStatus) error {
	op, err := protocol.NewResponse(protocol.SGIP_REPORT_REP, seq, status)
	if err != nil {
		return err
	}

	return s.conn.Write(op)
}
