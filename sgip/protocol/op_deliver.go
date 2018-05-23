package protocol

import (
	"bytes"
	"fmt"
)

type Deliver struct {
	*Header

	// body
	UserNumber     *OctetString
	SPNumber       *OctetString
	TP_pid         uint8
	TP_udhi        uint8
	MessageCoding  uint8
	MessageLength  uint32
	MessageContent *OctetString
	Reserve        *OctetString
}

func NewDeliver(
	seq [3]uint32,

	userNumber, spNumber string,
	TP_pid, TP_udhi, messageCoding int,
	messageContent []byte,
) (*Deliver, error) {
	op := &Deliver{}

	msgLen := uint32(len(messageContent))

	// body
	var length uint32

	op.UserNumber = &OctetString{Data: []byte(userNumber), FixedLen: 21}
	length = length + 21

	op.SPNumber = &OctetString{Data: []byte(spNumber), FixedLen: 21}
	length = length + 21

	op.TP_pid = uint8(TP_pid)
	length = length + 1

	op.TP_udhi = uint8(TP_udhi)
	length = length + 1

	op.MessageCoding = uint8(messageCoding)
	length = length + 1

	op.MessageLength = msgLen
	length = length + 4

	op.MessageContent = &OctetString{Data: []byte(messageContent), FixedLen: int(msgLen)}
	length = length + msgLen

	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	// header
	op.Header = &Header{}
	op.Length = length + 4 + 4 + 12
	op.CmdId = SGIP_DELIVER
	op.Sequence = seq

	return op, nil
}

func ParseDeliver(hdr *Header, data []byte) (*Deliver, error) {
	op := &Deliver{}
	op.Header = hdr

	p := 0
	op.UserNumber = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.SPNumber = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.TP_pid = data[p]
	p = p + 1

	op.TP_udhi = data[p]
	p = p + 1

	op.MessageCoding = data[p]
	p = p + 1

	op.MessageLength = unpackUi32(data[p : p+4])
	p = p + 4

	msgLen := int(op.MessageLength)
	op.MessageContent = &OctetString{Data: data[p : p+msgLen], FixedLen: msgLen}
	p = p + msgLen

	op.Reserve = &OctetString{Data: data[p : p+8], FixedLen: 8}
	p = p + 8

	return op, nil
}

func (p *Deliver) Serialize() []byte {
	b := p.Header.Serialize()

	b = append(b, p.UserNumber.Byte()...)
	b = append(b, p.SPNumber.Byte()...)
	b = append(b, packUi8(p.TP_pid)...)
	b = append(b, packUi8(p.TP_udhi)...)
	b = append(b, packUi8(p.MessageCoding)...)
	b = append(b, packUi32(p.MessageLength)...)
	b = append(b, p.MessageContent.Byte()...)
	b = append(b, p.Reserve.Byte()...)

	return b
}

func (p *Deliver) Ok() bool {
	return true
}

func (p *Deliver) String() string {
	var b bytes.Buffer
	b.WriteString(p.Header.String())

	fmt.Fprintln(&b, "--- Deliver ---")
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "SPNumber: ", p.SPNumber)
	fmt.Fprintln(&b, "TP_pid: ", p.TP_pid)
	fmt.Fprintln(&b, "TP_udhi: ", p.TP_udhi)
	fmt.Fprintln(&b, "MessageCoding: ", p.MessageCoding)
	fmt.Fprintln(&b, "MessageLength: ", p.MessageLength)
	fmt.Fprintln(&b, "MessageContent: ", p.MessageContent) // TODO convert encoding

	return b.String()
}
