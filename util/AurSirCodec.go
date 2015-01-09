package util

import (
	"encoding/json"
	"io/ioutil"
	"github.com/vmihailenco/msgpack"
	"github.com/ugorji/go/codec"
	"bytes"
)


type AurSirCodec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	DecodeFile(string, interface {}) error
}


func GetCodec(codec string) AurSirCodec {
	switch codec {
	case "JSON":
		return codecJson{}
	case "MSGPACK":
		return codecMsgpack{}
	default:
		return nil
	}

}


type codecJson struct{}

func (codecJson) Encode(i interface{}) ([]byte, error) {

	enc, err := json.Marshal(i)

	return enc, err
}

func (codecJson) Decode(b []byte, t interface{}) error {
	return json.Unmarshal(b, t)
}
func (codecJson) DecodeFile(filename string , t interface{}) error {
	src, err := ioutil.ReadFile(filename)
	if err!= nil {
		return err
	}
	return json.Unmarshal(src, t)
}



type codecMsgpack struct{}

func (codecMsgpack) Encode(i interface{}) ([]byte, error) {

	enc, err := msgpack.Marshal(i)

	return enc, err
}

func (codecMsgpack) Decode(b []byte, t interface{}) error {
	var h codec.MsgpackHandle
	dec := codec.NewDecoder(bytes.NewReader(b),&h)
	return dec.Decode(&t)
}
func (codecMsgpack) DecodeFile(filename string , t interface{}) error {
	src, err := ioutil.ReadFile(filename)
	if err!= nil {
		return err
	}
	return msgpack.Unmarshal(src, t)
}


