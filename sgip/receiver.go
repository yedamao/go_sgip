package sgip

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

type Handler interface {

	// bind
	OnBind(login_type uint8, name, password string) protocol.RespStatus

	// deliver
	OnDeliver(userNumber, spNumber string, TP_pid, TP_udhi,
		messageCoding uint8, messageContent []byte) protocol.RespStatus

	// report
	OnReport(seq [3]uint32, reportType uint8, userNumber string,
		state, errorCode uint8) protocol.RespStatus
}

// Receiver is a server
// listening on addr,
// accept SMGW connection
// and start sgip protocol Session
type Receiver struct {
	wg       sync.WaitGroup
	listener net.Listener

	addr  string
	count int // accept worker count
	done  chan struct{}

	handler Handler

	// debug flag
	debug bool
}

func NewReceiver(addr string, count int, handler Handler, debug bool) (*Receiver, error) {
	if handler == nil {
		return nil, errors.New("Reciver: nil handler")
	}

	r := &Receiver{
		addr:    addr,
		count:   count,
		done:    make(chan struct{}),
		handler: handler,
		debug:   debug,
	}
	err := r.setup()

	return r, err
}

func (r *Receiver) worker(id int) {
	defer r.wg.Done()
	log.Println("worker ", id, " running..")

	for {

		conn, err := r.listener.Accept()
		if err != nil {
			select {
			case <-r.done:
				log.Println("worker ", id, ": Server closed")
				return
			default:
			}

			log.Println("Accept error: ", err)
			continue
		}

		log.Println("worker ", id, ": Accept ", conn)

		// block
		startSession(conn, r)
	}
}

func (r *Receiver) setup() (err error) {
	if r.listener, err = net.Listen("tcp", r.addr); err != nil {
		return err
	}

	log.Println("Receiver listen on: ", r.listener.Addr())
	return nil
}

func (r *Receiver) Run() {

	for i := 0; i < r.count; i++ {
		r.wg.Add(1)
		go r.worker(i)
	}

	r.wg.Wait()
}

func (r *Receiver) Stop() {
	close(r.done)
	r.listener.Close()
	log.Println("Server stopped...")
}

// 代表与运营商的一条会话连接
type Session struct {
	// SMG -> SP 的连接
	conn     conn.Conn
	receiver *Receiver
}

// 开启会话
// block
func startSession(connection net.Conn, recv *Receiver) {
	s := &Session{
		conn:     *conn.NewConn(connection),
		receiver: recv,
	}

	s.start()
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

func (s *Session) start() {
	// 关闭会话
	defer s.close()

	for {

		select {
		case <-s.receiver.done:
			log.Println("Session: Server closed")
			return
		default:
		}

		s.conn.SetDeadline(time.Now().Add(1e9))
		op, err := s.conn.Read()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Println("Receiver: read error exit. ", err)
			return
		}

		if s.receiver.debug { // if debug mode, print op pkg
			log.Println(op)
		}

		switch op.GetHeader().CmdId {
		case protocol.SGIP_BIND:
			bind, ok := op.(*protocol.Bind)
			if !ok {
				log.Println("Receiver: bind type assert error")
				return
			}

			stat := s.receiver.handler.OnBind(
				bind.Type, bind.Name.String(), bind.Password.String(),
			)

			s.bindResp(op.GetHeader().Sequence, stat)

		case protocol.SGIP_DELIVER:
			deliver, ok := op.(*protocol.Deliver)
			if !ok {
				log.Println("Receiver: deliver type assert error")
				return
			}

			stat := s.receiver.handler.OnDeliver(
				deliver.UserNumber.String(), deliver.SPNumber.String(),
				deliver.TP_pid, deliver.TP_udhi, deliver.MessageCoding,
				deliver.MessageContent.Byte(),
			)

			s.deliverResp(op.GetHeader().Sequence, stat)

		case protocol.SGIP_REPORT:
			report, ok := op.(*protocol.Report)
			if !ok {
				log.Println("Receiver: report type assert error")
				return
			}

			stat := s.receiver.handler.OnReport(
				report.SubmitSequence, report.ReportType, report.UserNumber.String(),
				report.State, report.ErrorCode,
			)

			s.reportResp(op.GetHeader().Sequence, stat)

		case protocol.SGIP_UNBIND:
			s.unbindResp(op.GetHeader().Sequence)
			return

		default:
			log.Printf("Receiver: Unknow Operation CmdId: 0x%x\n", op.GetHeader().CmdId)
			return
		}
	}
}

// 关闭会话
func (s *Session) close() {
	s.conn.Close()
}
