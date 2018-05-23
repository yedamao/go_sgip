package protocol

import (
	"testing"
)

func TestBind(t *testing.T) {
	op, err := NewBind([...]uint32{0, 0, 0}, 1, "testusername", "testpwd")
	if err != nil {
		t.Error(err)
	}

	raw := op.Serialize()

	parsedHdr, err := ParseHeader(raw[:20])
	if err != nil {
		t.Error(err)
	}

	bind := op
	if parsedHdr.Length != bind.Length ||
		parsedHdr.CmdId != bind.CmdId {

		t.Log(bind.Header)
		t.Log(parsedHdr)
		t.Error("header not equal")
	}

	parsedBind, err := ParseBind(parsedHdr, raw[20:])
	if err != nil {
		t.Error(err)
	}

	if parsedBind.Type != bind.Type ||
		parsedBind.Name.String() != bind.Name.String() ||
		parsedBind.Password.String() != bind.Password.String() {

		t.Log(bind)
		t.Log(parsedBind)
		t.Error("bind not equal")
	}
}
