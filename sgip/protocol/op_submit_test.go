package protocol

import (
	"testing"
)

func TestSubmit(t *testing.T) {
	submit, err := NewSubmit(
		[...]uint32{0, 0, 0},
		"10690090",
		"17600537300",
		[]string{"17600000000", "17611111111"},
		"12345",
		"",
		0,
		"",
		"",
		0,
		0,
		0,
		"yymmddhhmmsstnnp",
		"yymmddhhmmsstnnp",
		0,
		0,
		0,
		0,
		0,
		[]byte("test msg"),
	)
	if err != nil {
		t.Error(err)
	}

	raw := submit.Serialize()
	parsedOp, err := ParseOperation(raw)
	if err != nil {
		t.Error(err)
	}

	parsedSubmit := parsedOp.(*Submit)

	if submit.Length != parsedSubmit.Length ||
		submit.SPNumber.String() != parsedSubmit.SPNumber.String() ||
		submit.MessageLength != parsedSubmit.MessageLength ||
		submit.MessageContent.String() != parsedSubmit.MessageContent.String() {

		t.Error("submit not equal")
	}
}
