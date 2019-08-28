package sgip

import (
	"strconv"
	"testing"

	"github.com/yedamao/go_sgip/sgip/protocol"
	"github.com/yedamao/go_sgip/sgip/sgiptest"
)

func TestRunTransmitter(t *testing.T) {

	host := "127.0.0.1"
	port := 1234

	// setup
	addr := host + ":" + strconv.Itoa(port)
	server, err := sgiptest.NewServer(addr)
	if err != nil {
		t.Fatal("New mock server: ", err)
	}

	go server.Run()

	var client *Transmitter
	// run testing
	t.Run("NewTransmitter", func(t *testing.T) {
		tx, err := NewTransmitter(host, port, "000", "000", "fakename", "1234")
		if err != nil {
			t.Fatal(err)
		}
		client = tx
	})

	t.Run("Submit", func(t *testing.T) {
		_, err := client.Submit("1069000000", []string{"17600000000"}, "no", 0, protocol.GBK, []byte("TestSubmitMessage"))
		if err != nil {
			t.Fatal(err)
		}

		assertMockServerResp(t, client, protocol.SGIP_SUBMIT_REP)
	})

	t.Run("UnBind", func(t *testing.T) {
		if err := client.Unbind(); err != nil {
			t.Fatal(err)
		}
		assertMockServerResp(t, client, protocol.SGIP_UNBIND_REP)
	})

	// teardown
	server.Stop()
}

func assertMockServerResp(t *testing.T, client *Transmitter, wantCMD uint32) {
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
