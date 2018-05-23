package sgip

import (
	"io"
	"testing"

	"github.com/yedamao/go_sgip/sgip/conn"
	"github.com/yedamao/go_sgip/sgip/sgiptest"
)

// mock sp client
// send bind/unbind/deliver/report
type ReceiverClient struct {
	conn.Conn
}

func newReceiver() (*Receiver, error) {

	return NewReceiver(":8008", 2, &sgiptest.MockHandler{}, false)
}

func TestRun(t *testing.T) {
	receiver, err := newReceiver()
	if err != nil {
		t.Error(err)
	}

	go receiver.Run()

	// new client
	client, err := sgiptest.NewOperatorClient()
	if err != nil {
		t.Error(err)
	}

	err = client.Bind("testname", "testpwd")
	if err != nil {
		t.Error(err)
	}

	err = client.DeliverOnly(
		"8617600000000", "10690000000", 0, 0, 0,
		[]byte("test msg..."),
	)
	if err != nil {
		t.Error(err)
	}
	err = client.ReportOnly([...]uint32{0, 0, 0}, "8617600000000", 0, 0)
	if err != nil {
		t.Error(err)
	}

	resp, err := client.Read()
	if err != nil {
		t.Error(err)
	}
	if !resp.Ok() {
		t.Error("Resp not ok")
	}

	resp, err = client.Read()
	if err != nil {
		t.Error(err)
	}
	if !resp.Ok() {
		t.Error("Resp not ok")
	}

	err = client.Unbind()
	if err != nil {
		t.Error(err)
	}

	// conn should closed
	_, err = client.Read()
	if err != io.EOF {
		t.Error(err)
	}

	receiver.Stop()
}

func BenchmarkRun(b *testing.B) {

	receiver, err := newReceiver()
	if err != nil {
		b.Error(err)
	}

	go receiver.Run()

	// new client
	client, err := sgiptest.NewOperatorClient()
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		err = client.DeliverOnly(
			"8617600000000", "10690000000", 0, 0, 0,
			[]byte("test msg..."),
		)
		if err != nil {
			b.Error(err)
		}
		resp, err := client.Read()
		if err != nil {
			b.Error(err)
		}
		if !resp.Ok() {
			b.Error("Resp not ok")
		}
	}

	err = client.Unbind()
	if err != nil {
		b.Error(err)
	}

	// conn should closed
	_, err = client.Read()
	if err != io.EOF {
		b.Error(err)
	}

	receiver.Stop()
}
