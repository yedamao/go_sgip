package protocol

import (
	"bytes"
	"fmt"
)

type Header struct {
	// Message Length 消息的总长度(字节)
	Length uint32
	// Command ID 命令ID
	CmdId uint32
	// // Sequence Number 序列号
	Sequence [3]uint32
}

func (p *Header) GetHeader() *Header {
	return p
}

func (p *Header) Serialize() []byte {
	b := packUi32(p.Length)
	b = append(b, packUi32(p.CmdId)...)
	b = append(b, packUi32(p.Sequence[0])...)
	b = append(b, packUi32(p.Sequence[1])...)
	b = append(b, packUi32(p.Sequence[2])...)

	return b
}

func (p *Header) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- Header ---")
	fmt.Fprintln(&b, "Length: ", p.Length)
	fmt.Fprintf(&b, "CmdId: 0x%x\n", p.CmdId)
	fmt.Fprintln(&b, "Sequence: ", p.Sequence)

	return b.String()
}

func (p *Header) Parse(data []byte) *Header {

	p.Length = unpackUi32(data[:4])
	p.CmdId = unpackUi32(data[4:8])

	p.Sequence[0] = unpackUi32(data[8:12])
	p.Sequence[1] = unpackUi32(data[12:16])
	p.Sequence[2] = unpackUi32(data[16:20])

	return p
}

func ParseHeader(data []byte) (*Header, error) {

	h := &Header{}
	h.Parse(data)

	return h, nil
}
