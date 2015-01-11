package dock


type Outgoing interface {
	Activate(id string) error
	Send(msgtype int64, codec string,msg []byte) (err error)
}

