package sgiptest

import (
	"github.com/yedamao/go_sgip/sgip/protocol"
)

type MockHandler struct{}

func (h *MockHandler) OnBind(login_type uint8, name, password string) protocol.RespStatus {
	// Do nothing

	return protocol.STAT_OK
}

func (h *MockHandler) OnDeliver(
	userNumber, spNumber string, TP_pid, TP_udhi,
	messageCoding uint8, messageContent []byte,
) protocol.RespStatus {
	// Do nothing

	return protocol.STAT_OK
}

func (h *MockHandler) OnReport(
	seq [3]uint32, reportType uint8, userNumber string,
	state, errorCode uint8,
) protocol.RespStatus {
	// Do nothing

	return protocol.STAT_OK
}
