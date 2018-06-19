package conn

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"

	"github.com/yedamao/go_sgip/sgip/errors"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

// Conn is a sgip connection can read/write protocol Operation
type Conn struct {
	net.Conn
	r *bufio.Reader
	w *bufio.Writer
}

// new a sgip Conn
func NewConn(fd net.Conn) *Conn {
	return &Conn{
		Conn: fd,
		r:    bufio.NewReader(fd),
		w:    bufio.NewWriter(fd),
	}
}

func (c *Conn) Read() (protocol.Operation, error) {
	l := make([]byte, 4)
	_, err := io.ReadFull(c.r, l)
	if err != nil {
		return nil, err

	}

	length := binary.BigEndian.Uint32(l) - 4
	if length > protocol.MAX_OP_SIZE {
		return nil, errors.SgipSizeErr
	}

	data := make([]byte, length)

	_, err = io.ReadFull(c.r, data)
	if err != nil {
		return nil, err
	}

	pkt := append(l, data...)

	op, err := protocol.ParseOperation(pkt)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (c *Conn) Write(op protocol.Operation) error {
	_, err := c.Conn.Write(op.Serialize())

	return err
}

func (c *Conn) Close() {
	c.Conn.Close()
}
