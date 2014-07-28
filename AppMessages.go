package AurSir4Go


import (

	"encoding/json"

	"errors"

)

type AppMessage struct {
	MsgType int64 //Command type of the package

	MsgCodec string //the codec used to serialize the message

	Msg	*[]byte //the encoded message as byte array
}

func (appMsg AppMessage) Decode() (AurSirMessage, error){


	codec := getCodec(appMsg.MsgCodec)

	switch appMsg.MsgType {

	case DOCK:
		var m AurSirDockMessage
		err := codec.decode(appMsg.Msg,&m)
		return m, err

	case DOCKED:
		var m AurSirDockedMessage
		err := codec.decode(appMsg.Msg,&m)
		return m, err

	case LEAVE:
		var m AurSirLeaveMessage
	err := codec.decode(appMsg.Msg,&m)
return m, err
	case REQUEST:
		var m AurSirRequest
err := codec.decode(appMsg.Msg,&m)
return m, err
	case RESULT:
		var m AurSirResult
err := codec.decode(appMsg.Msg,&m)
return m, err
	case ADD_EXPORT:
		var m AurSirAddExportMessage
err := codec.decode(appMsg.Msg,&m)
return m, err
	case EXPORT_ADDED:
		var m AurSirExportAddedMessage
err := codec.decode(appMsg.Msg,&m)
return m, err
	case ADD_IMPORT:
		var m AurSirAddImportMessage
err := codec.decode(appMsg.Msg,&m)
return m, err

	case IMPORT_ADDED:
		var m AurSirImportAddedMessage
err := codec.decode(appMsg.Msg,&m)
return m, err
case IMPORT_UPDATED:
		var m AurSirImportUpdatedMessage
err := codec.decode(appMsg.Msg,&m)
return m, err

	case LISTEN:
		var m AurSirListenMessage
err := codec.decode(appMsg.Msg,&m)
return m, err
	case STOP_LISTEN:
		var m AurSirStopListenMessage
err := codec.decode(appMsg.Msg,&m)
return m, err
	case UPDATE_EXPORT:
		var m AurSirUpdateExportMessage
err := codec.decode(appMsg.Msg,&m)
return m, err
	case UPDATE_IMPORT:
		var m AurSirUpdateImportMessage
err := codec.decode(appMsg.Msg,&m)
return m, err

	default:
		return nil, errors.New("Unknown Message")

	}
}

type codec interface{
	encode(interface {}) (*[]byte,error)
	decode(*[]byte,interface {}) error
}

type codecJson struct {}

func (codecJson) encode(i interface {})(*[]byte,error){

	enc, err := json.Marshal(i)

	return &enc, err
}

func (codecJson) decode(b *[]byte,t interface {})error{
	return json.Unmarshal(*b,t)
}

func getCodec(codec string) codec{
	switch codec {
	case "JSON":
		return codecJson{}
	default:
		return nil
	}

}

func (appMsg *AppMessage) Encode(msg AurSirMessage, codec string) (error){

  	appMsg.MsgType = msg2cmd(msg)
	appMsg.MsgCodec = codec
	c := getCodec(codec)
	var err error
	appMsg.Msg, err = c.encode(msg)
	return err

}

//cmd2msg returns a integer cmd for a given msg type
func msg2cmd(msg AurSirMessage) int64 {

	switch (msg).(type) {

	case AurSirDockMessage:
		return DOCK
	case AurSirDockedMessage:
		return DOCKED
	case AurSirLeaveMessage:
		return LEAVE
	case AurSirRequest:
		return REQUEST
	case AurSirResult:
		return RESULT
	case AurSirAddExportMessage:
		return ADD_EXPORT
	case AurSirExportAddedMessage:
		return EXPORT_ADDED
	case AurSirAddImportMessage:
		return ADD_IMPORT
	case AurSirImportAddedMessage:
		return IMPORT_ADDED
	case AurSirImportUpdatedMessage:
		return IMPORT_UPDATED
	case AurSirListenMessage:
		return LISTEN
	case AurSirStopListenMessage:
		return STOP_LISTEN
	case AurSirUpdateExportMessage:
		return UPDATE_EXPORT
	case AurSirUpdateImportMessage:
		return UPDATE_IMPORT

	}
	return -1
}
