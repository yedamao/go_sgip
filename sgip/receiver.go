package sgip

import (
	"errors"
	"log"
	"net"
	"sync"

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

		session := NewSession(conn, r.handler, r.done)
		// block
		session.Run()
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
