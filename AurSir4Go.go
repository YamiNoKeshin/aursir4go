package AurSir4Go

import (
	"github.com/pebbe/zmq4"
	"log"
	"strconv"
	"net"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

type AurSirInterface struct {

	AppName string //Name of the application

	UUID string //the applications uuid

	in, out chan AurSirMessage //channels for communication between front and backend

	quit *bool // a quit flag

	port int64 // the incoming port of the interface

	connected *bool // a connected flag

	exports *map[string] *ExportedAppKey

	exportsSem chan struct {}

	imports *map[string] *ImportedAppKey

    resultChans *map[string] chan Result
}

//NewInterface creates a new AurSir interface and returns a pointer to it
func NewInterface(AppName string) *AurSirInterface {

	//create interface
	var iface AurSirInterface

	//set the app name
	iface.AppName = AppName

	//generate a UUID
	iface.UUID = generateUuid()

	//initalize channels for communication with the backend
	iface.in = make(chan AurSirMessage,10)
	iface.out = make(chan AurSirMessage,10)

	//initialize the quit flag
	q := false
	iface.quit = &q

	//get a free network port
	iface.port = getRandomPort()

	//init the connection flag
	conn := false
	iface.connected = &conn


	exports := map[string]*ExportedAppKey{}
	iface.exports = &exports
	iface.exportsSem = make(chan struct {},1)
	iface.exportsSem <- struct {}{}

	imports := map[string]*ImportedAppKey{}
	iface.imports = &imports

	resChans := map[string]chan Result{}
	iface.resultChans = &resChans

	go iface.backend()

	iface.connect()

	return &iface

}

//Close shuts down the AurSir interface
func (iface *AurSirInterface) Close(){
	log.Println("Closing out channel")
	close(iface.out)
	log.Println("out channel closed")

	for !*iface.quit{
		time.Sleep(10*time.Millisecond)
	}
}

func (iface *AurSirInterface) AddExport(key AppKey,tags[]string) *ExportedAppKey{
	<- iface.exportsSem
	var ak ExportedAppKey
	ak.iface = iface
	ak.key = key
	ak.tags = tags

	expReq := AurSirAddExportMessage{key,tags}

	iface.out <- expReq

	expRep := <- iface.in

	expMsg , ok := expRep.(AurSirExportAddedMessage)
	if !ok {panic("insane runtime!!!")}

	ak.exportId = expMsg.ExportId

	ak.Request = make(chan Request)

	(*iface.exports)[expMsg.ExportId] = &ak
	log.Println(*iface.exports)

	iface.exportsSem <-struct{}{}
	return &ak
}

func (iface *AurSirInterface) AddImport(key AppKey,tags[]string) *ImportedAppKey{
	var ak ImportedAppKey
	ak.iface = iface
	ak.key = key
	ak.tags = tags

	ak.listenChan = make(chan Result)

	impReq := AurSirAddImportMessage{key,tags}

	iface.out <- impReq

	impRep := <- iface.in

	impMsg , ok := impRep.(AurSirImportAddedMessage)
	if !ok {panic("insane runtime!!!")}

	ak.importId = impMsg.ImportId
	ak.Connected = impMsg.Exported
	(*iface.imports)[impMsg.ImportId] = &ak
	return &ak
}

//registerResultChan stores a request uuid for ONE2.. or appkey.functioname@tag for MANY2... calls together with a channel so when Results comes in later, it can be routed to the
// appropriate channel.
func (iface *AurSirInterface) registerResultChan(resUuid string, rc chan Result){
	(*iface.resultChans)[resUuid] = rc
}


//connect initializes the connection to the runtime by sending a DOCK message
func (iface *AurSirInterface) connect(){
	iface.out <- AurSirDockMessage{iface.AppName}
}

func (iface *AurSirInterface) backend(){

	go iface.listener()

	go iface.sender()

}

func (iface *AurSirInterface) sender() {

	log.Println("Opening ougoing ZeroMQ Socket")

	skt, err := zmq4.NewSocket(zmq4.DEALER)

	if err != nil{
		panic("Could not open ZeroMQ socket")
	}

	defer skt.Close()

	skt.SetIdentity(iface.UUID)

	skt.Connect("tcp://localhost:5555")

	log.Println("Outgoing ZeroMQ Socket open")

	for msg := range iface.out {

		var appmsg AppMessage

		appmsg.Encode(msg,"JSON")
		log.Println("Sending",appmsg)

		if appmsg.MsgType == DOCK {
			skt.SendMessage([]string{strconv.FormatInt(appmsg.MsgType,10),appmsg.MsgCodec,string(appmsg.Msg),strconv.FormatInt(iface.port,10)},0)
		} else {
		skt.SendMessage([]string{strconv.FormatInt(appmsg.MsgType,10),appmsg.MsgCodec,string(appmsg.Msg)},0)
		}
	}
	log.Println("Sending LEAVE")
	var lmsg AppMessage

	lmsg.Encode(AurSirLeaveMessage{},"JSON")

	skt.SendMessage([]string{strconv.FormatInt(lmsg.MsgType,10),lmsg.MsgCodec,string(lmsg.Msg)},0)
	*iface.quit = true
}

func (iface *AurSirInterface) listener() {

	log.Println("Opening Incoming ZeroMQ Socket")

	skt, err := zmq4.NewSocket(zmq4.ROUTER)

	if err != nil{
		panic("Could not open ZeroMQ socket")
	}

	defer skt.Close()

	skt.Bind("tcp://*:"+strconv.FormatInt(iface.port,10))

	log.Println("Incoming ZeroMQ Socket open")

	skt.SetRcvtimeo(100)

	for !*iface.quit {

		msg, err := skt.RecvMessage(0)

		if err == nil {
			log.Println("Got incoming message")
			iface.processMsg(msg)
		}

	}

	log.Println("Closing incoming channel")

}

func (iface *AurSirInterface) processMsg(message []string){

	msgType, err := strconv.ParseInt(message[1],10,64)

	if err != nil{
		return
	}

	switch msgType{

	case DOCKED:
		*iface.connected = true

	case IMPORT_UPDATED:
		msg := AppMessage{msgType,message[2],[]byte(message[3])}
		asmsg, err := msg.Decode()
		if err == nil {
			iumsg,ok := asmsg.(AurSirImportUpdatedMessage)
			if ok {
				(*iface.imports)[iumsg.ImportId].Connected = iumsg.Exported
			}

		}

	case REQUEST:
		msg := AppMessage{msgType,message[2],[]byte(message[3])}
		asmsg, err := msg.Decode()

		if err == nil {
			reqmsg,ok := asmsg.(AurSirRequest)

			if ok {
				//this could and should be done more elgently
				<- iface.exportsSem
				for _,exp := range *iface.exports {
					log.Println(exp.key.ApplicationKeyName)
					if exp.key.ApplicationKeyName == reqmsg.AppKeyName{

						exp.Request <- Request{reqmsg.Request,reqmsg.FunctionName,reqmsg.Uuid,reqmsg.CallType}

					}
				}
				iface.exportsSem <-struct{}{}
			}

		}
case RESULT:
		msg := AppMessage{msgType,message[2],[]byte(message[3])}
		asmsg, err := msg.Decode()
		if err == nil {
			resmsg,ok := asmsg.(AurSirResult)

			if ok {

				var resId string

				if resmsg.CallType == ONE2MANY  || resmsg.CallType == ONE2ONE {
					resId = resmsg.Uuid
				}else{
					resId = resmsg.AppKeyName+"."+resmsg.FunctionName
				}


				rc, f := (*iface.resultChans)[resId]

				log.Println(resId)
				//log.Println(rc)
				if f {
					rc <- Result{resmsg.Result,resmsg.FunctionName,resmsg.Uuid,resmsg.CallType}
				}

			}

		}

	default:
		msg := AppMessage{msgType,message[2],[]byte(message[3])}
		asmsg, err := msg.Decode()

		if err == nil {
			iface.in <- asmsg
		}
	}
}


func (iface *AurSirInterface) Connected() bool{
	con := *iface.connected
	return con
}

func getRandomPort()  int64{
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err!=nil{
		panic("Could not find a free port")
	}
	defer l.Close()
	return int64(l.Addr().(*net.TCPAddr).Port)
}

func generateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}


/*

func (r AurSirRequest) ApplicationKeyName() string {
	var AppKey struct {
		ApplicationKeyName string
		}
	json.Unmarshal([]byte(r.RequestJSON), &AppKey)
	return AppKey.ApplicationKeyName
}

func (r AurSirRequest) FunctionName() string {
	var Fun struct {
		Function string
	}
	json.Unmarshal([]byte(r.RequestJSON), &Fun)
	return Fun.Function
}
*/


/*
func (r AurSirResult) ApplicationKeyName() string {
	var AppKey struct {
		ApplicationKeyName string
		}
	json.Unmarshal([]byte(r.ResultJSON), &AppKey)
	return AppKey.ApplicationKeyName
}

func (r AurSirResult) FunctionName() string {
	var Fun struct {
		Function string
	}
	json.Unmarshal([]byte(r.ResultJSON), &Fun)
	return Fun.Function
}
*/


