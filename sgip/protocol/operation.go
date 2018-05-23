package protocol

import (
	"errors"
	"fmt"
)

type Operation interface {
	//Header returns the Operation header, decoded. Header fields
	// can be updated before reserialzing .
	GetHeader() *Header

	// SerializeTo encodes Operation to it's binary form,
	// include the header and body
	Serialize() []byte

	// check resp result
	Ok() bool

	// String
	String() string
}

func ParseOperation(data []byte) (Operation, error) {
	if len(data) < 20 {
		return nil, errors.New("Invalide data length")
	}

	header, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}

	if int(header.Length) != len(data) {
		return nil, errors.New("Invalide data length")
	}

	var n Operation

	switch header.CmdId {
	case SGIP_UNBIND:
		n, err = ParseUnbind(header, data[20:])
	case SGIP_UNBIND_REP:
		n, err = ParseUnbindResp(header, data[20:])

	case SGIP_BIND:
		n, err = ParseBind(header, data[20:])
	case SGIP_SUBMIT:
		n, err = ParseSubmit(header, data[20:])
	case SGIP_DELIVER:
		n, err = ParseDeliver(header, data[20:])
	case SGIP_REPORT:
		n, err = ParseReport(header, data[20:])

	case SGIP_BIND_REP, SGIP_SUBMIT_REP, SGIP_DELIVER_REP, SGIP_REPORT_REP:
		n, err = ParseResponse(header, data[20:])

	default:
		err = fmt.Errorf("Unknow Operation CmdId: 0x%x", header.CmdId)
	}

	return n, err
}
