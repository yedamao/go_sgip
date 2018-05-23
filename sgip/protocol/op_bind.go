package protocol

import (
	"bytes"
	"fmt"
)

type Bind struct {
	*Header

	Type     uint8        // 登录类型 1 sp -> SMG, 2 SMG -> SP
	Name     *OctetString // 登录名
	Password *OctetString // 密码
	Reserve  *OctetString // 保留
}

func NewBind(seq [3]uint32, login_type uint8, name, password string) (*Bind, error) {
	op := &Bind{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 12 // header

	op.Type = login_type
	length = length + 1

	op.Name = &OctetString{Data: []byte(name), FixedLen: 16}
	length = length + 16

	op.Password = &OctetString{Data: []byte(password), FixedLen: 16}
	length = length + 16

	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	op.Length = length
	op.CmdId = SGIP_BIND
	op.Sequence = seq

	return op, nil
}

func ParseBind(hdr *Header, data []byte) (*Bind, error) {
	p := 0
	op := &Bind{}
	op.Header = hdr

	op.Type = data[p]
	p = p + 1

	op.Name = &OctetString{Data: data[p : p+16], FixedLen: 16}
	p = p + 16

	op.Password = &OctetString{Data: data[p : p+16], FixedLen: 16}
	p = p + 16

	op.Reserve = &OctetString{Data: data[p : p+8], FixedLen: 8}
	p = p + 8

	return op, nil
}

func (p *Bind) Serialize() []byte {
	b := p.Header.Serialize()

	b = append(b, packUi8(p.Type)...)
	b = append(b, p.Name.Byte()...)
	b = append(b, p.Password.Byte()...)
	b = append(b, p.Reserve.Byte()...)

	return b
}

func (p *Bind) Ok() bool {
	return true
}

func (p *Bind) String() string {
	var b bytes.Buffer
	b.WriteString(p.Header.String())

	fmt.Fprintln(&b, "--- Bind ---")
	fmt.Fprintln(&b, "Login Type: ", p.Type)
	fmt.Fprintln(&b, "Login Name: ", p.Name)
	fmt.Fprintln(&b, "Login Password: ", p.Password)

	return b.String()
}
