package client

import (
	"net"
	"strconv"
	"sync"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/errors"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

// Client Common Client
type Client struct {
	conn.Conn

	mu       sync.Mutex
	nodeId   uint32
	CorpId   string
	sequence uint32
}

func (c *Client) NewSeqNum() [3]uint32 {
	defer c.mu.Unlock()

	c.mu.Lock()
	c.sequence++

	return [3]uint32{c.nodeId, protocol.TimeStamp(), c.sequence}
}

func (c *Client) Setup(areaCode, corpId string) error {

	nodeId, err := protocol.NodeId(areaCode, corpId)
	if err != nil {
		return err
	}
	c.nodeId = nodeId

	c.CorpId = corpId

	return nil
}

func (c *Client) Connect(host string, port int) error {
	connection, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	c.Conn = *conn.NewConn(connection)

	return nil
}

func (c *Client) Bind(name, password string, login_type uint8) error {
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

func (c *Client) Unbind() error {
	op, err := protocol.NewUnbind(c.NewSeqNum())
	if err != nil {
		return err
	}

	return c.Write(op)
}

func (c *Client) UnbindResp(sequence [3]uint32) error {
	op, err := protocol.NewUnbindResp(sequence)
	if err != nil {
		return err
	}

	return c.Write(op)
}
