package internal

import (
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
	"github.com/joernweissenborn/aursir4go/calltypes"
)

type Incomingprocessor struct {
	Docked chan bool
	AddExport chan string
	AddImport chan string
	ExportedChans map[string]chan bool
	ResultChans map[string]chan messages.Result
	RequestChans map[string]chan messages.Request

}

func InitIncomingProcessor() *Incomingprocessor{
	var i Incomingprocessor
	i.Docked = make(chan bool,1)
	i.Docked <- false

	i.AddExport = make(chan string)
	i.AddImport = make(chan string)
	i.ResultChans = map[string]chan messages.Result{}
	i.RequestChans = map[string]chan messages.Request{}
	i.ExportedChans = map[string]chan bool{}
	return &i
}

//registerResultChan stores a request uuid for ONE2.. or appkey.functioname@tag for MANY2... calls together with a channel so when Results comes in later, it can be routed to the
// appropriate channel.
func (i *Incomingprocessor) RegisterResultChan(resUuid string, rc chan messages.Result) {
	i.ResultChans[resUuid] = rc
}
func (i *Incomingprocessor) RegisterRequestChan(expid string, rc chan messages.Request) {
	i.RequestChans[expid] = rc
}

func (i *Incomingprocessor) ProcessMsg(msgType int64, codec string, message []byte) {



	decoder := util.GetCodec(codec)
	if decoder == nil{
		return
	}
	switch msgType {

	case messages.DOCKED:
		var m messages.DockedMessage
		decoder.Decode(message,&m)
		<- i.Docked
	i.Docked <- m.Ok

	case messages.IMPORT_UPDATED:
		var m messages.ImportUpdatedMessage
		decoder.Decode(message,&m)
		c := i.ExportedChans[m.ImportId]
		<- c
	c <- m.Exported
	case messages.EXPORT_ADDED:
		var m messages.ExportAddedMessage
		decoder.Decode(message,&m)
		 i.AddExport <- m.ExportId
	case messages.IMPORT_ADDED:
		var m messages.ImportAddedMessage
		decoder.Decode(message,&m)
		i.ExportedChans[m.ImportId] = make(chan bool, 1)
		c := i.ExportedChans[m.ImportId]
		i.AddImport <- m.ImportId
	c <- m.Exported

	case messages.REQUEST:
		var m messages.Request
		decoder.Decode(message,&m)
		c := i.RequestChans[m.ExportId]
		if c != nil {
			c<-m
		}
	case messages.RESULT:
		var m messages.Result
		decoder.Decode(message,&m)

		resId := ""

		if m.CallType == calltypes.ONE2MANY || m.CallType == calltypes.ONE2ONE {
			resId = m.Uuid
		} else {
			resId = m.AppKeyName + "." + m.FunctionName
		}

		c, f := i.ResultChans[resId]

		if f {
			c<-m
		}
	}

}

