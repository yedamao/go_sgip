package encoding

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type GB18030 []byte

// Encode to GB18030
func (s GB18030) Encode() []byte {
	e := simplifiedchinese.GB18030
	es, _, err := transform.Bytes(e.NewEncoder(), s)
	if err != nil {
		return s
	}

	return es
}

// Decode from GB18030
func (s GB18030) Decode() []byte {
	e := simplifiedchinese.GB18030
	es, _, err := transform.Bytes(e.NewDecoder(), s)
	if err != nil {
		return s
	}

	return es
}

// UCS2 text codec.
type UCS2 []byte

// Encode to UCS2.
func (s UCS2) Encode() []byte {
	e := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	es, _, err := transform.Bytes(e.NewEncoder(), s)
	if err != nil {
		return s

	}
	return es

}

// Decode from UCS2.
func (s UCS2) Decode() []byte {
	e := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	es, _, err := transform.Bytes(e.NewDecoder(), s)
	if err != nil {
		return s

	}
	return es

}

func GBK2UCS2(msg []byte) []byte {
	// GBK -> UTF-8 -> UCS2
	return UCS2(GB18030(msg).Decode()).Encode()

}

func UTF82GBK(msg []byte) []byte {
	return GB18030(msg).Encode()

}
