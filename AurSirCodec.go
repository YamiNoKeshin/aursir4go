package aursir4go

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"github.com/vmihailenco/msgpack"
)


type AurSirCodec interface {
	Encode(interface{}) (*[]byte, error)
	Decode(*[]byte, interface{}) error
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


func (asr AurSirRequest) Decode(target interface{}) error {

	codec := GetCodec(asr.Codec)
	if codec == nil {
		return errors.New("Unknown codec "+asr.Codec)
	}

	if asr.IsFile {
		return codec.DecodeFile(string(asr.Request),&target)
	}

	return codec.Decode(&asr.Request, &target)
}



type codecJson struct{}

func (codecJson) Encode(i interface{}) (*[]byte, error) {

	enc, err := json.Marshal(i)

	return &enc, err
}

func (codecJson) Decode(b *[]byte, t interface{}) error {
	return json.Unmarshal(*b, t)
}
func (codecJson) DecodeFile(filename string , t interface{}) error {
	src, err := ioutil.ReadFile(filename)
	if err!= nil {
		return err
	}
	return json.Unmarshal(src, t)
}


func (appMsg *AppMessage) Encode(msg AurSirMessage, codec string) error {

	appMsg.MsgType = msg2cmd(msg)
	appMsg.MsgCodec = codec
	c := GetCodec(codec)
	var err error
	appMsg.Msg, err = c.Encode(msg)
	return err

}

type codecMsgpack struct{}

func (codecMsgpack) Encode(i interface{}) (*[]byte, error) {

	enc, err := msgpack.Marshal(i)

	return &enc, err
}

func (codecMsgpack) Decode(b *[]byte, t interface{}) error {
	return msgpack.Unmarshal(*b, t)
}
func (codecMsgpack) DecodeFile(filename string , t interface{}) error {
	src, err := ioutil.ReadFile(filename)
	if err!= nil {
		return err
	}
	return msgpack.Unmarshal(src, t)
}


