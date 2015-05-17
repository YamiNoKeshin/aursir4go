package aursir4go

import (
	"errors"
	"github.com/joernweissenborn/aursir4go/appkey"
	"github.com/joernweissenborn/aursir4go/calltypes"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
	"time"
)

//An ImportedAppKey represents an applications imports and is used to call and listen to functions and to create
// callchains
type ImportedAppKey struct {
	iface                 *AurSirInterface
	key                   appkey.AppKey
	tags                  []string
	importId              string
	connected             chan bool
	listenFuns            []string
	listenChan            chan messages.Result
	persistenceStrategies map[string]string
}

//GetId returns the currently registered tags for the import.
func (iak *ImportedAppKey) GetId() string {
	return iak.importId
}

//Tags returns the currently registered tags for the import.
func (iak *ImportedAppKey) Tags() []string {
	return iak.tags
}

//Name returns the name of the imports ApplicationKey.
func (iak *ImportedAppKey) Name() string {
	return iak.key.ApplicationKeyName
}

//Exported returns the imports exported state.
func (iak *ImportedAppKey) Exported() (exported bool) {
	exported = <-iak.connected
	iak.connected <- exported
	return
}

//Listen to functions registers the import for listening to this function. Use Listen to get Results for this function.
func (iak *ImportedAppKey) ListenToFunction(FunctionName string) {
	listenid := iak.key.ApplicationKeyName + "." + FunctionName
	msg, _ := util.GetCodec(iak.iface.codec).Encode(messages.ListenMessage{iak.importId, FunctionName})
	iak.iface.outgoing.Send(messages.LISTEN, iak.iface.codec, msg)
	iak.listenFuns = append(iak.listenFuns, listenid)
}

//GetListenChannel returns the imports listen channel.
func (iak *ImportedAppKey) GetListenChannel() (ListenChannel chan messages.Result) {
	ListenChannel = iak.listenChan
	return
}

//Listen listens for results on listened functions. If no listen functions have been added, it returns an empty result Result.
func (iak *ImportedAppKey) Listen() messages.Result {
	return <-iak.listenChan
}

//Call requests the function specified by FunctionName as ONE2ONE request and returns a channel to get the result.
func (iak *ImportedAppKey) Call(FunctionName string, Arguments interface{}) (chan messages.Result, error) {
	return iak.callFunction(FunctionName, Arguments, calltypes.ONE2ONE)
}

//CallAll requests the function specified by FunctionName as ONE2MANY request. Get the results via Listen() or the
// ListenChannel.
func (iak *ImportedAppKey) CallAll(FunctionName string, Arguments interface{}) (err error) {
	_, err = iak.callFunction(FunctionName, Arguments, calltypes.ONE2MANY)
	return
}

//Trigger requests the function specified by FunctionName as MANY2ONE request. Get the results via Listen() or the
// ListenChannel.
func (iak *ImportedAppKey) Trigger(FunctionName string, Arguments interface{}) (err error) {
	_, err = iak.callFunction(FunctionName, Arguments, calltypes.MANY2ONE)
	return
}

//TriggerAll requests the function specified by FunctionName as MANY2MANY request. Get the results via Listen() or the
// ListenChannel.
func (iak *ImportedAppKey) TriggerAll(FunctionName string, Arguments interface{}) (err error) {
	_, err = iak.callFunction(FunctionName, Arguments, calltypes.MANY2MANY)
	return
}

//Call functions calls the function specified by FunctionName and returns a channel to get the result. This channel we
// be nil on Many2... call types! You need to use Listen() in this case.
func (iak *ImportedAppKey) CallFunction(FunctionName string, Arguments interface{}, CallType int64) (chan messages.Result, error) {
	return iak.callFunction(FunctionName, Arguments, CallType)
}

//UpdateTags sets the imports tags while overriding the old and registers the new tagset at the runtime. If you want to
// add a tag, use AddTag.
func (iak *ImportedAppKey) UpdateTags(NewTags []string) {
	iak.tags = NewTags
	msg, _ := util.GetCodec(iak.iface.codec).Encode(messages.UpdateImportMessage{iak.importId, iak.tags})
	iak.iface.outgoing.Send(messages.UPDATE_IMPORT, iak.iface.codec, msg)
}

//AddTag adds a tag to the imports tags and registers the new tagset at the runtime. If you want to set a new tagset,
// use UpdateTags
func (iak *ImportedAppKey) AddTag(Tag string) {
	iak.UpdateTags(append(iak.tags, Tag))
}

func (iak *ImportedAppKey) callFunction(FunctionName string, Arguments interface{}, CallType int64) (chan messages.Result, error) {

	if CallType > 3 {
		return nil, errors.New("Invalid calltype")
	}

	codec := util.GetCodec("JSON")
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
	iak.iface.outgoing.Send(messages.REQUEST, iak.iface.codec, msg)

	var resChan chan messages.Result
	if CallType == calltypes.ONE2ONE || CallType == calltypes.ONE2MANY {
		resChan = make(chan messages.Result)
		iak.iface.incomingprocessor.RegisterResultChan(reqUuid, resChan)
	}

	return resChan, nil
}
