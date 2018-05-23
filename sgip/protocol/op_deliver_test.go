package protocol

import (
	"testing"
)

func TestDeliver(t *testing.T) {
	deliver, err := NewDeliver(
		[...]uint32{0, 0, 0},
		"17600111111",
		"1069999999999",
		0,
		0,
		0,
		[]byte("hello deliver test msg"),
	)
	if err != nil {
		t.Error(err)
	}

	raw := deliver.Serialize()
	parsedOp, err := ParseOperation(raw)
	if err != nil {
		t.Error(err)
	}

	parsedDeliver := parsedOp.(*Deliver)

	if deliver.Length != parsedDeliver.Length ||
		deliver.UserNumber.String() != parsedDeliver.UserNumber.String() ||
		deliver.SPNumber.String() != parsedDeliver.SPNumber.String() ||
		deliver.MessageLength != parsedDeliver.MessageLength {
		t.Error("deliver not equal")
	}
}
