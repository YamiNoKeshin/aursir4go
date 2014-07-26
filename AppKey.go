package AurSir4Go

import (
	"errors"
	"log"
)



type AppKey struct {

	ApplicationKeyName string

	Functions []Function
}

type Function struct {

	Name string

	Input []Data

	Output []Data

}

type Data struct {

	Name string

	Type int

}


type ImportedAppKey struct{
	iface *AurSirInterface
	key AppKey
	tags []string
	importId string
	Connected bool
	listenFuns []string
	listenChan chan Result
}

func (iak *ImportedAppKey) ListenToFunction(FunctionName string)  {


	iak.iface.out <- AurSirListenMessage{iak.importId,FunctionName}

	listenid := iak.key.ApplicationKeyName+"."+FunctionName


	iak.iface.registerResultChan(listenid,iak.listenChan)
	log.Println("wfhoooooooooooooooooooooo",listenid)

	iak.listenFuns = append(iak.listenFuns,listenid)


}

func (iak *ImportedAppKey) Listen() Result {
	if len(iak.listenFuns) == 0 {
		return Result{}
	}
	return <- iak.listenChan
}

func (iak *ImportedAppKey) CallFunction(FunctionName string, Arguments interface {}, CallType int64) (chan Result,error){

	if CallType >3 { return nil,errors.New("Invalid calltype")}
	log.Println(Arguments)
	codec := getCodec("JSON")
	args, err := codec.encode(Arguments)
	if err!=nil{return nil,err}

	reqUuid := generateUuid()

	iak.iface.out <- AurSirRequest{iak.key.ApplicationKeyName,FunctionName,CallType,iak.tags,reqUuid,"JSON",args}

	if CallType == ONE2ONE || CallType == ONE2MANY {
		resChan := make(chan Result)
		iak.iface.registerResultChan(reqUuid, resChan)

	return resChan,nil}
	return nil,nil
}


type ExportedAppKey struct{
	iface *AurSirInterface
	key AppKey
	tags []string
	exportId string
	Request chan Request
}

func (eak ExportedAppKey) Reply(req *Request, res interface {}) error {
	var aursirResult AurSirResult
	aursirResult.AppKeyName = eak.key.ApplicationKeyName
	aursirResult.Codec = "JSON"
	aursirResult.CallType = req.CallType
	aursirResult.Uuid = req.Uuid
	aursirResult.Tags = eak.tags
	aursirResult.FunctionName = req.Function
	codec:= getCodec("JSON")
	result, err := codec.encode(res)

	if err ==nil {
		aursirResult.Result = result
		eak.iface.out <- aursirResult
	}
	return err
}

type Request struct {
	req []byte
	Function string
	Uuid string
	CallType int64
}

func (r Request) Decode(target interface {}){
	codec:= getCodec("JSON")
	codec.decode(r.req,&target)
}

type Result struct {
	res []byte
	Function string
	Uuid string
	CallType int64
}
func (r Result) Decode(target interface {}){
	codec:= getCodec("JSON")
	codec.decode(r.res,&target)
}
