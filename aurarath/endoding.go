package aurarath

import "bytes"

const (
	BIN uint8 = iota
	JSON
)

var (
	CODECS []uint8 = []uint8{JSON}
)

func decode(s *bytes.Buffer, t interface{}, codec string) {}
func encode(interface{}) (b *bytes.Buffer) {
	return
}
