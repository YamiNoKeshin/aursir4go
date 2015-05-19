package aurarath

import (
	"bytes"
	"encoding/json"
)

const (
	BIN uint8 = iota
	JSON
)

var (
	CODECS []uint8 = []uint8{JSON}
)


func decode(s *bytes.Buffer, t interface{}, codec uint8) {
	switch codec{

	case JSON:
		jsondec := json.NewDecoder(s)
		jsondec.Decode(t)
	}
}
func encode(d interface{}) (b *bytes.Buffer) {
	enc, _ := json.Marshal(d)

	return bytes.NewBuffer(enc)
}
