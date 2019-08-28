package sgip

import (
	"net"
	"strconv"
	"sync"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/errors"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

// commonClient
type commonClient struct {
	conn.Conn

	Mu       sync.Mutex
	nodeId   uint32
	corpId   string
	sequence uint32
}

func (c *commonClient) NewSeqNum() [3]uint32 {
	defer c.Mu.Unlock()

	c.Mu.Lock()
	c.sequence++

	return [3]uint32{c.nodeId, protocol.TimeStamp(), c.sequence}
}

func (c *commonClient) setup(areaCode, corpId string) error {

	nodeId, err := protocol.NodeId(areaCode, corpId)
	if err != nil {
		return err
	}
	c.nodeId = nodeId

	c.corpId = corpId

	return nil
}

func (c *commonClient) connect(host string, port int) error {
	connection, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	c.Conn = *conn.NewConn(connection)

	return nil
}

func (c *commonClient) bind(name, password string, login_type uint8) error {
	op, err := protocol.NewBind(c.NewSeqNum(), login_type, name, password)
	if err = c.Write(op); err != nil {
		return err
	}

	// If BindResp NOT received in 5secs close connection

	// Read block
	var resp protocol.Operation
	if resp, err = c.Read(); err != nil {
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

func (c *commonClient) Unbind() error {
	op, err := protocol.NewUnbind(c.NewSeqNum())
	if err != nil {
		return err
	}

	return c.Write(op)
}

func (c *commonClient) UnbindResp(sequence [3]uint32) error {
	op, err := protocol.NewUnbindResp(sequence)
	if err != nil {
		return err
	}

	return c.Write(op)
}
