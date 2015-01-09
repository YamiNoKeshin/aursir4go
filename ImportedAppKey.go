package aursir4go

import (
	"errors"
	"time"
	"github.com/joernweissenborn/aursir4go/appkey"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
)
//An ImportedAppKey represents an applications imports and is used to call and listen to functions and to create
// callchains
type ImportedAppKey struct {
	iface      *AurSirInterface
	key        appkey.AppKey
	tags       []string
	importId   string
	Connected  bool
	listenFuns []string
	listenChan chan messages.Result
	persistenceStrategies map[string] string
}

//Tags returns the currently registered tags for the import.
func (iak *ImportedAppKey) Tags() []string {
	return iak.tags
}

//Name returns the name of the imports ApplicationKey.
func (iak *ImportedAppKey) Name() string {
	return iak.key.ApplicationKeyName
}

//Listen to functions registers the import for listening to this function. Use Listen to get Results for this function.
func (iak *ImportedAppKey) ListenToFunction(FunctionName string) {
	listenid := iak.key.ApplicationKeyName + "." + FunctionName
	iak.iface.incomingprocessor.RegisterResultChan(listenid, iak.listenChan)
	msg, _ := util.GetCodec(iak.iface.codec).Encode(messages.ListenMessage{iak.importId, FunctionName})
	iak.iface.outgoing.Send(messages.UPDATE_EXPORT,iak.iface.codec,msg)
	iak.listenFuns = append(iak.listenFuns, listenid)
}

//Listen listens for results on listened functions. If no listen functions have been added, it returns an empty result Result.
func (iak *ImportedAppKey) Listen() Result {
	if len(iak.listenFuns) == 0 {
		return Result{}
	}
	return <-iak.listenChan
}

//Call functions calls the function specified by FunctionName and returns a channel to get the result. This channel we
// be nil on Many2... call types! You need to use Listen() in this case.
func (iak *ImportedAppKey) CallFunction(FunctionName string, Arguments interface{}, CallType int64) (chan Result, error) {
	return iak.callFunction(FunctionName,Arguments,CallType, false)
}


//UpdateTags sets the imports tags while overriding the old and registers the new tagset at the runtime. If you want to
// add a tag, use AddTag.
func (iak *ImportedAppKey) UpdateTags(NewTags []string) {
	iak.tags = NewTags
	msg, _ := util.GetCodec(iak.iface.codec).Encode(messages.UpdateExportMessage{iak.importId, iak.tags})
	iak.iface.outgoing.Send(messages.UPDATE_IMPORT,iak.iface.codec,msg)	}

//AddTag adds a tag to the imports tags and registers the new tagset at the runtime. If you want to set a new tagset,
// use UpdateTags
func (iak *ImportedAppKey) AddTag(Tag string) {
	iak.UpdateTags(append(iak.tags,Tag))
}


func (iak *ImportedAppKey) callFunction(FunctionName string, Arguments interface{}, CallType int64, Persist bool) (chan Result, error) {

	if CallType > 3 {
		return nil, errors.New("Invalid calltype")
	}

	codec := GetCodec("JSON")
	if codec == nil {
		return nil, errors.New("unknown codec")
	}

	args, err := codec.Encode(Arguments)
	if err != nil {
		return nil, err
	}

	reqUuid := generateUuid()

	req := messages.Request{
		iak.key.ApplicationKeyName,
		FunctionName,
		CallType,
		iak.tags,
		reqUuid,
		iak.importId,
		"",
		time.Now(),
		"JSON",
		args,
		false,
		false,
	}

	msg, _ := util.GetCodec(iak.iface.codec).Encode(req)
	iak.iface.outgoing.Send(messages.UPDATE_IMPORT,iak.iface.codec,msg)

	var resChan chan messages.Result
	if CallType == ONE2ONE || CallType == ONE2MANY {
		resChan = make(chan messages.Result)
		iak.iface.incomingprocessor.RegisterResultChan(reqUuid, resChan)
	}

	return resChan, nil
}
