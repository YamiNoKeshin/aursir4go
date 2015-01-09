package aursir4go

import (
	"time"
	"github.com/joernweissenborn/aursir4go/appkey"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
)

//ExportedAppKey provides methods to get incoming requests and reply to them. It also other methods to set a persistence
// mode, which will override any Requested persistence settings.
type ExportedAppKey struct {
	iface    *AurSirInterface
	key      appkey.AppKey
	tags     []string
	exportId string
	Request  chan messages.Request
}

//Tags returns the exports current tags.
func (eak ExportedAppKey) Tags() []string {
	return eak.tags
}


//Reply encodes the result and sends it to the interface, where it is transmitted to the aursir runtime. Perssitence is
// either set by the Request or is overridden by calling e.g. SetLogging.
func (eak ExportedAppKey) Reply(Request *messages.Request, Result interface{}) error {
	return eak.reply(Request,Result,false,true)
}

//StreamingReply works like Reply, but sets the streaming flag on the Result, so importers can react on it. Streams are
// identified by their requests, so use the same Request to make consecutive stream. Set finished true to indicate that
// you finished sending data
func (eak ExportedAppKey) StreamingReply(Request *messages.Request, Result interface{}, Finished bool) error {
	return eak.reply(Request,Result,true,Finished)
}

//UpdateTags sets the exports tags while overriding the old and registers the new tagset at the runtime. If you want to
// add a tag, use AddTag.
func (eak *ExportedAppKey) UpdateTags(NewTags []string) {
	eak.tags = NewTags
	msg, _ := util.GetCodec(eak.iface.codec).Encode(messages.UpdateExportMessage{eak.exportId, eak.tags})
	eak.iface.outgoing.Send(messages.UPDATE_EXPORT,eak.iface.codec,msg)
}

//AddTag adds a tag to the exports tags and registers the new tagset at the runtime. If you want to set a new tagset,
// use UpdateTags
func (eak *ExportedAppKey) AddTag(Tag string) {
	eak.UpdateTags(append(eak.tags,Tag))
}



func (eak ExportedAppKey) reply(req *messages.Request, res interface{},stream, finished bool) error {
	var aursirResult messages.Result
	aursirResult.AppKeyName = eak.key.ApplicationKeyName
	aursirResult.Codec = "JSON"
	aursirResult.CallType = req.CallType
	aursirResult.Uuid = req.Uuid
	aursirResult.ImportId = req.ImportId
	aursirResult.ExportId = eak.exportId
	aursirResult.Tags = eak.tags
	aursirResult.FunctionName = req.FunctionName
	aursirResult.Timestamp = time.Now()
	aursirResult.Stream = stream
	aursirResult.StreamFinished = finished

	codec := util.GetCodec("JSON")
	var err error
	aursirResult.Result, err = codec.Encode(res)

	if err == nil {
		msg, _ := util.GetCodec(eak.iface.codec).Encode(aursirResult)
		eak.iface.outgoing.Send(messages.RESULT,eak.iface.codec,msg)		}
	return err
}

func (eak ExportedAppKey) Emit(FunctionName string, Result interface{},stream, finished bool) error {
	Request := new(messages.Request)
	Request.FunctionName = FunctionName
	return eak.reply(Request,Result,stream,finished)
}
