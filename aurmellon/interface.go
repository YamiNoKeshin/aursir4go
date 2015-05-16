package aurmellon

const (
	PROTOCOL_ID uint8 = 0
)

type Interface struct {

}

func (i Interface) ProtocolId() int8 {
	return PROTOCOL_ID
}
