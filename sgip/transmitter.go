package sgip

import (
	"bufio"
	"net"
	"strconv"
	"sync"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/errors"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

type Transmitter struct {
	Mu       sync.Mutex
	nodeId   uint32
	corpId   string
	sequence uint32

	conn.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (tx *Transmitter) NewSeqNum() [3]uint32 {
	defer tx.Mu.Unlock()

	tx.Mu.Lock()
	tx.sequence++

	return [3]uint32{tx.nodeId, protocol.TimeStamp(), tx.sequence}
}

func NewTransmitter(host string, port int, areaCode, corpId, name, password string) (*Transmitter, error) {
	tx := &Transmitter{}

	if err := tx.Setup(areaCode, corpId); err != nil {
		return nil, err
	}

	if err := tx.Connect(host, port); err != nil {
		return nil, err
	}

	if err := tx.Bind(name, password); err != nil {
		return nil, err
	}

	return tx, nil
}

func (tx *Transmitter) Setup(areaCode, corpId string) error {

	nodeId, err := protocol.NodeId(areaCode, corpId)
	if err != nil {
		return err
	}
	tx.nodeId = nodeId

	tx.corpId = corpId

	return nil
}

func (tx *Transmitter) Connect(host string, port int) error {
	connection, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	tx.Conn = conn.Conn{Conn: connection}

	return nil
}

func (tx *Transmitter) Bind(name, password string) error {
	op, err := protocol.NewBind(tx.NewSeqNum(), 1, name, password)
	if err = tx.Write(op); err != nil {
		return err
	}

	// If BindResp NOT received in 5secs close connection

	// Read block
	var resp protocol.Operation
	if resp, err = tx.Read(); err != nil {
		return err
	}

	if resp.GetHeader().CmdId != protocol.SGIP_BIND_REP {
		return errors.SgipBindRespErr
	}

	if !resp.Ok() {
		return errors.SgipBindErr
	}

	return nil
}

func (tx *Transmitter) Unbind() error {
	op, err := protocol.NewUnbind(tx.NewSeqNum())
	if err != nil {
		return err
	}

	return tx.Write(op)
}

func (tx *Transmitter) UnbindResp(sequence [3]uint32) error {
	op, err := protocol.NewUnbindResp(sequence)
	if err != nil {
		return err
	}

	return tx.Write(op)
}

func (tx *Transmitter) Submit(spNumber string, destId []string, serviceType string, TP_udhi, msgCoding int, msg []byte) ([3]uint32, error) {
	op, err := protocol.NewSubmit(
		tx.NewSeqNum(),
		spNumber, "000000000000000000000",
		destId, tx.corpId, serviceType, 1, "0", "0",
		0, 2, 0, "", "", 2, 0, TP_udhi, msgCoding, 0, msg,
	)

	if err != nil {
		return [3]uint32{}, err
	}

	if err := tx.Write(op); err != nil {
		return [3]uint32{}, err
	}

	return op.GetHeader().Sequence, nil
}

func (tx *Transmitter) bindCheck() {
	// TODO
}
