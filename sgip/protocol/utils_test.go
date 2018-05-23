package protocol

import (
	"testing"
)

func TestNodeId(t *testing.T) {
	tests := []struct {
		areaCode string
		corpId   string
		exp      uint32
	}{
		{"010", "92008", 3010092008},
		{"020", "92008", 3020092008},
		{"021", "92008", 3021092008},
		{"0311", "92008", 3031192008},
		{"0711", "92008", 3071192008},
	}

	for _, test := range tests {

		id, err := NodeId(test.areaCode, test.corpId)
		if err != nil {
			t.Error(err)
		}

		if id != test.exp {
			t.Error("not equal: expect ", test.exp, "real ", id)
		}
	}

}
