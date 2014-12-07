package aursir4go

import "time"

//ExportedAppKey provides methods to get incoming requests and reply to them. It also other methods to set a persistence
// mode, which will override any Requested persistence settings.
type ExportedAppKey struct {
	iface    *AurSirInterface
	key      AppKey
	tags     []string
	exportId string
	Request  chan AurSirRequest
	persistenceStrategies map[string] string
}

//Tags returns the exports current tags.
func (eak ExportedAppKey) Tags() []string {
	return eak.tags
}


//Reply encodes the result and sends it to the interface, where it is transmitted to the aursir runtime. Perssitence is
// either set by the Request or is overridden by calling e.g. SetLogging.
func (eak ExportedAppKey) Reply(Request *AurSirRequest, Result interface{}) error {
	return eak.reply(Request,Result,false,true)
}

//StreamingReply works like Reply, but sets the streaming flag on the Result, so importers can react on it. Streams are
// identified by their requests, so use the same Request to make consecutive stream. Set finished true to indicate that
// you finished sending data
func (eak ExportedAppKey) StreamingReply(Request *AurSirRequest, Result interface{}, Finished bool) error {
	return eak.reply(Request,Result,true,Finished)
}

//UpdateTags sets the exports tags while overriding the old and registers the new tagset at the runtime. If you want to
// add a tag, use AddTag.
func (eak *ExportedAppKey) UpdateTags(NewTags []string) {
	eak.tags = NewTags
	eak.iface.out <- AurSirUpdateExportMessage{eak.exportId, eak.tags}
}

//AddTag adds a tag to the exports tags and registers the new tagset at the runtime. If you want to set a new tagset,
// use UpdateTags
func (eak *ExportedAppKey) AddTag(Tag string) {
	eak.UpdateTags(append(eak.tags,Tag))
}

//SetLogging sets the persistence strategy for function calls of the specified function to "log". This overrides all previous
// persistence strategies!
func (eak *ExportedAppKey) SetLogging(FunctionName string){
	eak.persistenceStrategies[FunctionName] = "log"
}



func (eak ExportedAppKey) reply(req *AurSirRequest, res interface{},stream, finished bool) error {
	var aursirResult AurSirResult
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

	if strategy,f := eak.persistenceStrategies[req.FunctionName]; f {
		aursirResult.Persistent = true
		aursirResult.PersistenceStrategy = strategy
	}

	codec := GetCodec("JSON")
	var err error
	aursirResult.Result, err = codec.Encode(res)

	if err == nil {
		eak.iface.out <- aursirResult
	}
	return err
}

func (eak ExportedAppKey) Emit(FunctionName string, Result interface{},stream, finished bool) error {
	var aursirResult AurSirResult
	aursirResult.AppKeyName = eak.key.ApplicationKeyName
	aursirResult.Codec = "JSON"
	aursirResult.CallType = MANY2ONE
	aursirResult.Uuid = ""
	aursirResult.ExportId = eak.exportId
	aursirResult.Tags = eak.tags
	aursirResult.FunctionName = FunctionName
	aursirResult.Timestamp = time.Now()
	aursirResult.Stream = stream
	aursirResult.StreamFinished = finished



	codec := GetCodec("JSON")
	var err error
	aursirResult.Result, err = codec.Encode(Result)

	if err == nil {
		eak.iface.out <- aursirResult
	}
	return err
}
