package protocol

// Unbind operation
type Unbind struct {
	*Header
}

func NewUnbind(seq [3]uint32) (*Unbind, error) {
	op := &Unbind{Header: &Header{}}

	op.Length = 4 + 4 + 12
	op.CmdId = SGIP_UNBIND
	op.Sequence = seq

	return op, nil
}

func ParseUnbind(hdr *Header, data []byte) (*Unbind, error) {

	return &Unbind{Header: hdr}, nil
}

func (p *Unbind) Serialize() []byte {

	return p.Header.Serialize()
}

func (p *Unbind) Ok() bool {
	return true
}

func (p *Unbind) String() string {

	return p.Header.String()
}

// UnbindResp operation
type UnbindResp struct {
	*Header
}

func NewUnbindResp(seq [3]uint32) (*UnbindResp, error) {
	op := &UnbindResp{Header: &Header{}}

	op.Length = 4 + 4 + 12
	op.CmdId = SGIP_UNBIND_REP
	op.Sequence = seq

	return op, nil
}

func ParseUnbindResp(hdr *Header, data []byte) (*UnbindResp, error) {

	return &UnbindResp{Header: hdr}, nil
}

func (p *UnbindResp) Serialize() []byte {

	return p.Header.Serialize()
}

func (p *UnbindResp) Ok() bool {
	return true
}

func (p *UnbindResp) String() string {

	return p.Header.String()
}
