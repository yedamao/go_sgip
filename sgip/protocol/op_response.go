package protocol

import (
	"bytes"
	"fmt"
)

type RespStatus uint8

func (s *RespStatus) Byte() []byte {
	return []byte{uint8(*s)}
}

// Common Response. Bind, Submit, Deliver, Report
type Response struct {
	*Header
	Result  RespStatus
	Reserve *OctetString
}

func NewResponse(CmdId uint32, seq [3]uint32, result RespStatus) (*Response, error) {
	op := &Response{Header: &Header{}}

	var length uint32 = 4 + 4 + 12 // header len

	op.Result = result
	length = length + 1
	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	// header
	op.Length = length
	op.CmdId = CmdId
	op.Sequence = seq

	return op, nil
}

func ParseResponse(hdr *Header, data []byte) (*Response, error) {
	op := &Response{}
	op.Header = hdr

	p := 0

	op.Result = RespStatus(data[p])
	p = p + 1

	op.Reserve = &OctetString{Data: data[p : p+8], FixedLen: 8}
	p = p + 8

	return op, nil
}

func (p *Response) GetHeader() *Header {
	return p.Header
}

func (p *Response) Serialize() []byte {
	b := p.Header.Serialize()

	b = append(b, p.Result.Byte()...)
	b = append(b, p.Reserve.Byte()...)

	return b
}

func (p *Response) Ok() bool {
	if STAT_OK == p.Result {
		return true
	}
	return false
}

func (p *Response) String() string {

	var b bytes.Buffer
	b.WriteString(p.Header.String())

	fmt.Fprintln(&b, "--- Response ---")
	fmt.Fprintln(&b, "Result: ", p.Result)

	return b.String()
}
