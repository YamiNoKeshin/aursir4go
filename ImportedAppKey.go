package AurSir4Go

import (
	"errors"
)

type ImportedAppKey struct{
	iface *AurSirInterface
	key AppKey
	tags []string
	importId string
	Connected bool
	listenFuns []string
	listenChan chan Result
}

func (iak ImportedAppKey) Tags() []string{
	return iak.tags
}

func (iak *ImportedAppKey) ListenToFunction(FunctionName string)  {
	listenid := iak.key.ApplicationKeyName+"."+FunctionName
	iak.iface.registerResultChan(listenid,iak.listenChan)
	iak.iface.out <- AurSirListenMessage{iak.importId,FunctionName}
	iak.listenFuns = append(iak.listenFuns,listenid)
}

func (iak *ImportedAppKey) Listen() Result {
	if len(iak.listenFuns) == 0 {
		return Result{}
	}
	return <- iak.listenChan
}

func (iak *ImportedAppKey) CallFunction(FunctionName string, Arguments interface {}, CallType int64) (chan Result,error){

	if CallType >3 {
		return nil,errors.New("Invalid calltype")
	}

	codec := getCodec("JSON")
	args, err := codec.encode(Arguments)
	if err!=nil{
		return nil,err
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
		*args}

	var resChan chan Result
	if CallType == ONE2ONE || CallType == ONE2MANY {
		resChan = make(chan Result)
		iak.iface.registerResultChan(reqUuid, resChan)
	}

	return resChan,nil
}

func (iak *ImportedAppKey) UpdateTags(NewTags []string){
	iak.tags = NewTags
	iak.iface.out <- AurSirUpdateImportMessage{iak.importId,iak.tags}
}
