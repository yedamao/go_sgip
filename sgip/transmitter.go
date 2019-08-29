package sgip

import (
	. "github.com/yedamao/go_sgip/sgip/client"
	"github.com/yedamao/go_sgip/sgip/protocol"
)

type Transmitter struct {
	Client
}

func NewTransmitter(host string, port int, areaCode, corpId, name, password string) (*Transmitter, error) {
	tx := &Transmitter{}

	if err := tx.Setup(areaCode, corpId); err != nil {
		return nil, err
	}

	if err := tx.Connect(host, port); err != nil {
		return nil, err
	}

	if err := tx.Bind(name, password); err != nil {
		return nil, err
	}

	return tx, nil
}

func (tx *Transmitter) Bind(name, password string) error {
	// sp bind to smg type 1
	return tx.Client.Bind(name, password, 1)
}

func (tx *Transmitter) Submit(spNumber string, destId []string, serviceType string, TP_udhi, msgCoding int, msg []byte) ([3]uint32, error) {
	op, err := protocol.NewSubmit(
		tx.NewSeqNum(),
		spNumber, "000000000000000000000",
		destId, tx.Client.CorpId, serviceType, 1, "0", "0",
		0, 2, 0, "", "", 1, 0, TP_udhi, msgCoding, 0, msg,
	)

	if err != nil {
		return [3]uint32{}, err
	}

	if err := tx.Write(op); err != nil {
		return [3]uint32{}, err
	}

	return op.GetHeader().Sequence, nil
}
