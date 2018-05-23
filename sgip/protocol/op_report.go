package protocol

import (
	"bytes"
	"fmt"
)

type Report struct {
	*Header

	// body
	SubmitSequence [3]uint32
	ReportType     uint8
	UserNumber     *OctetString
	State          uint8
	ErrorCode      uint8
	Reserve        *OctetString
}

func NewReport(
	seq [3]uint32, // 本条消息的序号
	msgSeq [3]uint32, // 该命令所涉及的Submit或deliver命令的序列号

	reportType int,
	userNumber string,
	state, errorCode int,
) (*Report, error) {
	op := &Report{}

	// body
	var length uint32

	op.SubmitSequence = msgSeq
	length = length + 4 + 4 + 4

	op.ReportType = uint8(reportType)
	length = length + 1

	op.UserNumber = &OctetString{Data: []byte(userNumber), FixedLen: 21}
	length = length + 21

	op.State = uint8(state)
	length = length + 1

	op.ErrorCode = uint8(errorCode)
	length = length + 1

	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	// header
	op.Header = &Header{}
	op.Length = length + 4 + 4 + 12
	op.CmdId = SGIP_REPORT
	op.Sequence = seq

	return op, nil
}

func ParseReport(hdr *Header, data []byte) (*Report, error) {

	op := &Report{}
	op.Header = hdr

	p := 0

	op.SubmitSequence[0] = unpackUi32(data[p : p+4])
	p = p + 4
	op.SubmitSequence[1] = unpackUi32(data[p : p+4])
	p = p + 4
	op.SubmitSequence[2] = unpackUi32(data[p : p+4])
	p = p + 4

	op.ReportType = data[p]
	p = p + 1

	op.UserNumber = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.State = data[p]
	p = p + 1

	op.ErrorCode = data[p]
	p = p + 1

	op.Reserve = &OctetString{Data: data[p : p+8], FixedLen: 8}
	p = p + 8

	return op, nil
}

func (p *Report) Serialize() []byte {
	b := p.Header.Serialize()

	b = append(b, packUi32(p.SubmitSequence[0])...)
	b = append(b, packUi32(p.SubmitSequence[1])...)
	b = append(b, packUi32(p.SubmitSequence[2])...)
	b = append(b, packUi8(p.ReportType)...)
	b = append(b, p.UserNumber.Byte()...)
	b = append(b, packUi8(p.State)...)
	b = append(b, packUi8(p.ErrorCode)...)
	b = append(b, p.Reserve.Byte()...)

	return b
}

func (p *Report) Ok() bool {
	return true
}

func (p *Report) String() string {
	var b bytes.Buffer
	b.WriteString(p.Header.String())

	fmt.Fprintln(&b, "--- Report ---")
	fmt.Fprintln(&b, "SubmitSequence: ", p.SubmitSequence)
	fmt.Fprintln(&b, "ReportType :", p.ReportType)
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "State :", p.State)
	fmt.Fprintln(&b, "ErrorCode :", p.ErrorCode)

	return b.String()
}
