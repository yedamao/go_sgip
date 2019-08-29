package sgip

import (
	"strconv"
	"testing"

	"github.com/yedamao/go_sgip/sgip/errors"
	"github.com/yedamao/go_sgip/sgip/protocol"
	"github.com/yedamao/go_sgip/sgip/sgiptest"
)

var (
	mockReceiver *Receiver

	host = "localhost"
	port = 8008
)

// setup running receiver server
func setup() error {
	if mockReceiver != nil {
		return nil
	}
	addr := host + ":" + strconv.Itoa(port)

	receiver, err := NewReceiver(addr, 1, &sgiptest.MockHandler{}, true)
	if err != nil {
		return err
	}

	go receiver.Run()

	return nil
}

func teardown() {
	if mockReceiver != nil {
		mockReceiver.Stop()
	}
}

func TestRunReceiver(t *testing.T) {

	if err := setup(); err != nil {
		t.Fatal(err)
	}

	t.Run("NewSMGClient with wrong name/password", func(t *testing.T) {
		// new SMG client
		_, err := sgiptest.NewSMGClient(host, port, "000", "000", "fakename", "wrong password")
		if err != errors.SgipBindErr {
			t.Error("NewSMGClient should bind auth failed")
		}
	})

	var client *sgiptest.SMGClient

	t.Run("NewSMGClient normal", func(t *testing.T) {
		// new SMG client
		c, err := sgiptest.NewSMGClient(host, port, "000", "000", "fakename", "1234")
		if err != nil {
			t.Error("NewSMGClient :", err)
		}

		client = c
	})

	t.Run("Test Deliver", func(t *testing.T) {
		err := client.Deliver("17600000000", "106900000", 0, 0, protocol.ASCII, []byte("TestDeliver"))
		if err != nil {
			t.Fatal(err)
		}

		assertResponse(t, client, protocol.SGIP_DELIVER_REP)
	})

	t.Run("Test Report", func(t *testing.T) {
		err := client.Report([3]uint32{0, 0, 0}, 0, "17600000000", 0, 0)
		if err != nil {
			t.Fatal(err)
		}

		assertResponse(t, client, protocol.SGIP_REPORT_REP)
	})

	t.Run("Test Unbind", func(t *testing.T) {
		if err := client.Unbind(); err != nil {
			t.Fatal(err)
		}

		assertResponse(t, client, protocol.SGIP_UNBIND_REP)
	})

	teardown()
}

func assertResponse(t *testing.T, client *sgiptest.SMGClient, wantCMD uint32) {
	t.Helper()

	// read one response
	op, err := client.Read()
	if err != nil {
		t.Fatal(err)
	}
	if op.GetHeader().CmdId != wantCMD {
		t.Errorf("response not expect %s want 0x%0x", op, wantCMD)
	}
	if !op.Ok() {
		t.Error("reponse not ok")
	}
}
