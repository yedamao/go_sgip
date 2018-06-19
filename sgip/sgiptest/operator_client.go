package sgiptest

import (
	"net"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/errors"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

type OperatorClient struct {
	conn.Conn
}

func NewOperatorClient() (*OperatorClient, error) {
	c := &OperatorClient{}

	err := c.Connect("localhost:8008")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (tx *OperatorClient) Connect(addr string) error {
	connection, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	tx.Conn = *conn.NewConn(connection)

	return nil
}

func (c *OperatorClient) Bind(name, password string) error {
	op, err := protocol.NewBind([...]uint32{0, 0, 0}, 2, name, password)
	if err = c.Write(op); err != nil {
		return err
	}

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

func (c *OperatorClient) Unbind() error {
	op, err := protocol.NewUnbind([...]uint32{0, 0, 0})
	if err != nil {
		return err
	}
	err = c.Write(op)
	if err != nil {
		return err
	}

	resp, err := c.Read()
	if err != nil {
		return err
	}

	// checko op
	if resp.GetHeader().CmdId != protocol.SGIP_UNBIND_REP {
		return errors.SgipWrongCmdId
	}

	return err
}

// write Deliver msg not wait response
func (c *OperatorClient) DeliverOnly(
	userNumber, spNumber string,
	TP_pid, TP_udhi, messageCoding int,
	messageContent []byte,
) error {
	op, err := protocol.NewDeliver(
		[...]uint32{0, 0, 0}, userNumber, spNumber,
		TP_pid, TP_udhi, messageCoding, messageContent,
	)
	if err != nil {
		return err
	}

	return c.Write(op)
}

func (c *OperatorClient) ReportOnly(
	sendSeq [3]uint32, userNumber string,
	stat, errorCode int,
) error {

	op, err := protocol.NewReport(
		[...]uint32{0, 0, 0}, sendSeq,
		0, userNumber, stat, errorCode,
	)

	if err != nil {
		return err
	}

	return c.Write(op)
}
