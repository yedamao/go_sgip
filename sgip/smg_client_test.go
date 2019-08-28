package sgip

import (
	"github.com/yedamao/go_sgip/sgip/protocol"
)

// SMGClient 模拟SMG->SP
//
// 支持命令:
// Bind
// Unbind
// Deliver
// Report
type SMGClient struct {
	commonClient
}

func NewSMGClient(host string, port int, areaCode, corpId, name, password string) (*SMGClient, error) {
	sc := &SMGClient{}

	if err := sc.setup(areaCode, corpId); err != nil {
		return nil, err
	}

	if err := sc.connect(host, port); err != nil {
		return nil, err
	}

	if err := sc.Bind(name, password); err != nil {
		return nil, err
	}

	return sc, nil
}

func (sc *SMGClient) Bind(name, password string) error {
	// smg -> sp server
	return sc.bind(name, password, 2)
}

func (sc *SMGClient) Deliver(userNumber string, spNumber string, TP_pid int, TP_udhi int, msgCoding int, msg []byte) error {
	op, err := protocol.NewDeliver(
		sc.NewSeqNum(),
		userNumber, spNumber, TP_pid, TP_udhi, msgCoding, msg,
	)
	if err != nil {
		return err
	}

	return sc.Write(op)
}

func (sc *SMGClient) Report(msgSeq [3]uint32, reportType int, userNumber string, state int, errorCode int) error {
	op, err := protocol.NewReport(
		sc.NewSeqNum(),
		msgSeq, reportType, userNumber, state, errorCode,
	)
	if err != nil {
		return err
	}

	return sc.Write(op)
}
