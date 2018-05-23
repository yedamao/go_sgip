package protocol

import (
	"testing"
)

func TestResponse(t *testing.T) {
	resp, err := NewResponse(SGIP_BIND_REP, [...]uint32{0, 0, 0}, STAT_OK)
	if err != nil {
		t.Error(err)
	}

	raw := resp.Serialize()

	parsedOp, err := ParseOperation(raw)
	if err != nil {
		t.Error(err)
	}

	parsedResp := parsedOp.(*Response)

	if parsedResp.Result != resp.Result {
		t.Error("bindResp not equal")
	}
}
