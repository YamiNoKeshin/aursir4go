package aursir4go

import (
	"errors"
)

type ImportedAppKey struct {
	iface      *AurSirInterface
	key        AppKey
	tags       []string
	importId   string
	Connected  bool
	listenFuns []string
	listenChan chan Result
}

func (iak ImportedAppKey) Tags() []string {
	return iak.tags
}

func (iak *ImportedAppKey) ListenToFunction(FunctionName string) {
	listenid := iak.key.ApplicationKeyName + "." + FunctionName
	iak.iface.registerResultChan(listenid, iak.listenChan)
	iak.iface.out <- AurSirListenMessage{iak.importId, FunctionName}
	iak.listenFuns = append(iak.listenFuns, listenid)
}

func (iak *ImportedAppKey) Listen() Result {
	if len(iak.listenFuns) == 0 {
		return Result{}
	}
	return <-iak.listenChan
}

func (iak *ImportedAppKey) CallFunction(FunctionName string, Arguments interface{}, CallType int64) (chan Result, error) {
	return iak.callFunction(FunctionName,Arguments,CallType, false)
}

func (iak *ImportedAppKey) PersistentCallFunction(FunctionName string, Arguments interface{}, CallType int64) (chan Result, error) {
	return iak.callFunction(FunctionName,Arguments,CallType, true)
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

	iak.iface.out <- AurSirRequest{
		iak.key.ApplicationKeyName,
		FunctionName,
		CallType,
		iak.tags,
		reqUuid,
		iak.importId,
		"JSON",
		false,
		Persist,
		*args}

	var resChan chan Result
	if CallType == ONE2ONE || CallType == ONE2MANY {
		resChan = make(chan Result)
		iak.iface.registerResultChan(reqUuid, resChan)
	}

	return resChan, nil
}

func (iak *ImportedAppKey) UpdateTags(NewTags []string) {
	iak.tags = NewTags
	iak.iface.out <- AurSirUpdateImportMessage{iak.importId, iak.tags}
}

func (iak *ImportedAppKey) NewCallChain(OriginFunctionName string, Arguments interface{}, OriginCallType int64) (CallChain, error) {
	if OriginCallType > 3 {
		return CallChain{}, errors.New("Invalid calltype")
	}

	codec := GetCodec("JSON")
	args, err := codec.Encode(Arguments)
	if err != nil {
		return CallChain{}, err
	}

	cc := createCallChain(iak.iface)
	cc.setOrigin(iak.key.ApplicationKeyName, OriginFunctionName, "JSON", args, iak.Tags(), OriginCallType, iak.importId)
	return cc, nil

}

func (iak *ImportedAppKey) FinalizeCallChain(FunctionName string, ArgumentMap map[string]string, CallType int64, CallChain CallChain) (chan Result, error) {
	if CallType > 3 {
		return nil, errors.New("Invalid calltype")
	}

	reqUuid := generateUuid()
	fcc := ChainCall{
		iak.key.ApplicationKeyName,
		FunctionName,
		ArgumentMap,
		CallType,
		iak.Tags(),
		reqUuid}
	CallChain.finalImportId = iak.importId
	CallChain.finalCall = fcc
	err := CallChain.Finalize()

	if err != nil {
		return nil, err
	}

	var resChan chan Result
	if CallType == ONE2ONE || CallType == ONE2MANY {
		resChan = make(chan Result)
		iak.iface.registerResultChan(reqUuid, resChan)
	}

	return resChan, nil
}
