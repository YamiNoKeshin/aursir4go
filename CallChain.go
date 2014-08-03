package aursir4go

import (
	"errors"
	"log"
)

type ChainCall struct {
	AppKeyName   string
	FunctionName string
	ArgumentMap  map[string]string
	CallType     int64
	Tags         []string
	ChainCallId  string
}

type CallChain struct {
	iface         *AurSirInterface
	originRequest AurSirRequest
	chain         []ChainCall
	finalImportId string
	finalCall     ChainCall
}

func createCallChain(iface *AurSirInterface) CallChain {
	var cc CallChain
	cc.iface = iface
	cc.chain = []ChainCall{}
	return cc
}

func (cc *CallChain) setOrigin(oan, ofn, oc string, oa *[]byte, ot []string, oct int64, oiid string) {
	cc.originRequest.AppKeyName = oan
	cc.originRequest.FunctionName = ofn
	cc.originRequest.Codec = oc
	cc.originRequest.Request = *oa
	cc.originRequest.Tags = ot
	cc.originRequest.CallType = oct
	cc.originRequest.Uuid = generateUuid()
	cc.originRequest.ImportId = oiid
}

func (cc *CallChain) AddCall(AppKeyName, FunctionName string, ParameterMap map[string]string, CallType int64, Tags []string) {
	cc.chain = append(cc.chain, ChainCall{AppKeyName, FunctionName, ParameterMap, CallType, Tags, ""})
}

func (cc *CallChain) Finalize() error {
	cc.iface.out <- AurSirCallChain{
		cc.originRequest,
		cc.chain,
		cc.finalImportId,
		cc.finalCall}
	res, ok := (<-cc.iface.in).(AurSirCallChainAddedMessage)
	if !ok {
		panic("Runtime error")
	}
	log.Println(res.CallChainOk)
	if !res.CallChainOk {
		return errors.New("Indoable CallChain")
	}
	return nil
}
