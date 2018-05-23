package protocol

import (
	"bytes"
	"fmt"
)

type Submit struct {
	*Header

	// body
	SPNumber         *OctetString
	ChargeNumber     *OctetString
	UserCount        uint8
	UserNumber       []*OctetString // 接收该短消息的手机号，该字段重复UserCount指定的次数，手机号码前加“86”国别标志
	CorpId           *OctetString
	ServiceType      *OctetString
	FeeType          uint8
	FeeValue         *OctetString
	GivenValue       *OctetString
	AgentFlag        uint8
	MorelatetoMTFlag uint8
	Priority         uint8
	ExpireTime       *OctetString
	ScheduleTime     *OctetString
	ReportFlag       uint8
	TP_pid           uint8
	TP_udhi          uint8
	MessageCoding    uint8
	MessageType      uint8
	MessageLength    uint32
	MessageContent   *OctetString
	Reserve          *OctetString
}

func NewSubmit(
	seq [3]uint32,

	spNumber, chargeNumber string,
	userNumber []string, corpId, serviceType string, feeType int,
	feeValue, givenValue string,
	agentFlag, morelatetoMTFlag, priority int,
	expireTime, scheduleTime string,
	reportFlag, TP_pid, TP_udhi, messageCoding, messageType int,
	messageContent []byte,

) (*Submit, error) {

	op := &Submit{}

	// body
	var length uint32

	op.SPNumber = &OctetString{Data: []byte(spNumber), FixedLen: 21}
	length = length + 21

	op.ChargeNumber = &OctetString{Data: []byte(chargeNumber), FixedLen: 21}
	length = length + 21

	op.UserCount = uint8(len(userNumber))
	length = length + 1

	for _, number := range userNumber {
		op.UserNumber = append(op.UserNumber,
			&OctetString{Data: []byte(number), FixedLen: 21})
		length = length + 21
	}

	op.CorpId = &OctetString{Data: []byte(corpId), FixedLen: 5}
	length = length + 5

	op.ServiceType = &OctetString{Data: []byte(serviceType), FixedLen: 10}
	length = length + 10

	op.FeeType = uint8(feeType)
	length = length + 1

	op.FeeValue = &OctetString{Data: []byte(feeValue), FixedLen: 6}
	length = length + 6
	op.GivenValue = &OctetString{Data: []byte(givenValue), FixedLen: 6}
	length = length + 6

	op.AgentFlag = uint8(agentFlag)
	length = length + 1
	op.MorelatetoMTFlag = uint8(morelatetoMTFlag)
	length = length + 1
	op.Priority = uint8(priority)
	length = length + 1

	op.ExpireTime = &OctetString{Data: []byte(expireTime), FixedLen: 16}
	length = length + 16
	op.ScheduleTime = &OctetString{Data: []byte(scheduleTime), FixedLen: 16}
	length = length + 16

	op.ReportFlag = uint8(reportFlag)
	length = length + 1
	op.TP_pid = uint8(TP_pid)
	length = length + 1
	op.TP_udhi = uint8(TP_udhi)
	length = length + 1
	op.MessageCoding = uint8(messageCoding)
	length = length + 1
	op.MessageType = uint8(messageType)
	length = length + 1

	msgLen := uint32(len(messageContent))
	op.MessageLength = msgLen
	length = length + 4

	op.MessageContent = &OctetString{Data: []byte(messageContent), FixedLen: int(msgLen)}
	length = length + msgLen

	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	// header
	op.Header = &Header{}
	op.Length = length + 4 + 4 + 12
	op.CmdId = SGIP_SUBMIT
	op.Sequence = seq

	return op, nil
}

func ParseSubmit(hdr *Header, data []byte) (*Submit, error) {
	op := &Submit{}
	op.Header = hdr

	p := 0
	op.SPNumber = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.ChargeNumber = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.UserCount = data[p]
	p = p + 1

	// parse UserNumber
	for i := 0; i < int(op.UserCount); i++ {
		nubmer := &OctetString{Data: data[p : p+21]}
		p = p + 21
		op.UserNumber = append(op.UserNumber, nubmer)
	}

	op.CorpId = &OctetString{Data: data[p : p+5], FixedLen: 5}
	p = p + 5

	op.ServiceType = &OctetString{Data: data[p : p+10], FixedLen: 10}
	p = p + 10

	op.FeeType = data[p]
	p = p + 1

	op.FeeValue = &OctetString{Data: data[p : p+6], FixedLen: 6}
	p = p + 6

	op.GivenValue = &OctetString{Data: data[p : p+6], FixedLen: 6}
	p = p + 6

	op.AgentFlag = data[p]
	p = p + 1

	op.MorelatetoMTFlag = data[p]
	p = p + 1

	op.Priority = data[p]
	p = p + 1

	op.ExpireTime = &OctetString{Data: data[p : p+16], FixedLen: 16}
	p = p + 16

	op.ScheduleTime = &OctetString{Data: data[p : p+16], FixedLen: 16}
	p = p + 16

	op.ReportFlag = data[p]
	p = p + 1

	op.TP_pid = data[p]
	p = p + 1

	op.TP_udhi = data[p]
	p = p + 1

	op.MessageCoding = data[p]
	p = p + 1

	op.MessageType = data[p]
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

func (p *Submit) Serialize() []byte {
	b := p.Header.Serialize()

	b = append(b, p.SPNumber.Byte()...)
	b = append(b, p.ChargeNumber.Byte()...)
	b = append(b, packUi8(p.UserCount)...)

	// serialize UserNumber
	for i := 0; i < int(p.UserCount); i++ {
		b = append(b, p.UserNumber[i].Byte()...)
	}

	b = append(b, p.CorpId.Byte()...)
	b = append(b, p.ServiceType.Byte()...)
	b = append(b, packUi8(p.FeeType)...)
	b = append(b, p.FeeValue.Byte()...)
	b = append(b, p.GivenValue.Byte()...)
	b = append(b, packUi8(p.AgentFlag)...)
	b = append(b, packUi8(p.MorelatetoMTFlag)...)
	b = append(b, packUi8(p.Priority)...)
	b = append(b, p.ExpireTime.Byte()...)
	b = append(b, p.ScheduleTime.Byte()...)
	b = append(b, packUi8(p.ReportFlag)...)
	b = append(b, packUi8(p.TP_pid)...)
	b = append(b, packUi8(p.TP_udhi)...)
	b = append(b, packUi8(p.MessageCoding)...)
	b = append(b, packUi8(p.MessageType)...)
	b = append(b, packUi32(p.MessageLength)...)
	b = append(b, p.MessageContent.Byte()...)
	b = append(b, p.Reserve.Byte()...)

	return b
}

func (p *Submit) Ok() bool {
	return true
}

func (p *Submit) String() string {
	var b bytes.Buffer
	b.WriteString(p.Header.String())

	fmt.Fprintln(&b, "--- Submit ---")
	fmt.Fprintln(&b, "SPNumber: ", p.SPNumber)
	fmt.Fprintln(&b, "ChargeNumber: ", p.ChargeNumber)
	fmt.Fprintln(&b, "UserCount: ", p.UserCount)

	// print UserNubmer
	fmt.Fprintln(&b, "UserNumber: ")
	for i := 0; i < int(p.UserCount); i++ {
		fmt.Fprintln(&b, p.UserNumber[i].String())
	}

	fmt.Fprintln(&b, "CorpId: ", p.CorpId)
	fmt.Fprintln(&b, "ServiceType: ", p.ServiceType)
	fmt.Fprintln(&b, "FeeType: ", p.FeeType)
	fmt.Fprintln(&b, "FeeValue: ", p.FeeValue)
	fmt.Fprintln(&b, "GivenValue: ", p.GivenValue)
	fmt.Fprintln(&b, "AgentFlag: ", p.AgentFlag)
	fmt.Fprintln(&b, "MorelatetoMTFlag: ", p.MorelatetoMTFlag)
	fmt.Fprintln(&b, "Priority: ", p.Priority)
	fmt.Fprintln(&b, "ExpireTime: ", p.ExpireTime)
	fmt.Fprintln(&b, "ScheduleTime: ", p.ScheduleTime)
	fmt.Fprintln(&b, "ReportFlag: ", p.ReportFlag)
	fmt.Fprintln(&b, "TP_pid: ", p.TP_pid)
	fmt.Fprintln(&b, "TP_udhi: ", p.TP_udhi)
	fmt.Fprintln(&b, "MessageCoding: ", p.MessageCoding)
	fmt.Fprintln(&b, "MessageType: ", p.MessageType)
	fmt.Fprintln(&b, "MessageLength: ", p.MessageLength)
	fmt.Fprintln(&b, "MessageContent: ", p.MessageContent)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)

	return b.String()
}
