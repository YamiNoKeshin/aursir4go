package aursir4go

import (
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pebbe/zmq4"
	"log"
	"net"
	"strconv"
	"time"
)

//The AurSirInterface handles the runtime time connection. It holds methods to create im- and exports.
type AurSirInterface struct {
	AppName string //Name of the application

	UUID string //the applications uuid

	in, out chan AurSirMessage //channels for communication between front and backend

	quit *bool // a quit flag

	port int64 // the incoming port of the interface

	connected chan bool //a connected flag

	exports *map[string]*ExportedAppKey

	exportsSem chan struct{}

	imports *map[string]*ImportedAppKey

	resultChans *map[string]chan Result
}

//NewInterface creates and launches a new AurSirInterface and returns a pointer to it.
func NewInterface(AppName string) *AurSirInterface {

	//create interface
	var iface AurSirInterface

	//set the app name
	iface.AppName = AppName

	//generate a UUID
	iface.UUID = generateUuid()

	//initalize channels for communication with the backend
	iface.in = make(chan AurSirMessage, 10)
	iface.out = make(chan AurSirMessage, 10)

	//initialize the quit flag
	q := false
	iface.quit = &q

	//get a free network port
	iface.port = getRandomPort()

	//init the connection flag
	iface.connected = make(chan bool,1)
	//iface.connected <- false
	exports := map[string]*ExportedAppKey{}
	iface.exports = &exports
	iface.exportsSem = make(chan struct{}, 1)
	iface.exportsSem <- struct{}{}

	imports := map[string]*ImportedAppKey{}
	iface.imports = &imports

	resChans := map[string]chan Result{}
	iface.resultChans = &resChans

	go iface.backend()

	iface.connect()

	return &iface

}

//Close shuts down the AurSir interface
func (iface *AurSirInterface) Close() {
	log.Println("Closing out channel")
	close(iface.out)
	log.Println("out channel closed")

	for !*iface.quit {
		time.Sleep(10 * time.Millisecond)
	}
}
//AddExport adds the specified ApplicationKey and tags and registeres it at the runtime. It returns a pointer to an
// ExportedAppKey which can be userd to handle incoming requests
func (iface *AurSirInterface) AddExport(key AppKey, tags []string) *ExportedAppKey {
	for !iface.Connected(){}

	<-iface.exportsSem
	var ak ExportedAppKey
	ak.iface = iface
	ak.key = key
	ak.tags = tags
	ak.persistenceStrategies = map[string]string{}

	expReq := AurSirAddExportMessage{key, tags}

	iface.out <- expReq

	expRep := <-iface.in

	expMsg, ok := expRep.(AurSirExportAddedMessage)
	if !ok {
		panic("insane runtime!!!")
	}

	ak.exportId = expMsg.ExportId

	ak.Request = make(chan AurSirRequest)

	(*iface.exports)[expMsg.ExportId] = &ak
	log.Println(*iface.exports)

	iface.exportsSem <- struct{}{}
	return &ak
}

//AddImport adds the specified ApplicationKey and tags and registeres it at the runtime. It returns a pointer to an
// ImportedAppKey which can be used to request function calls or call chains and listening to functions
func (iface *AurSirInterface) AddImport(key AppKey, tags []string) *ImportedAppKey {

	for !iface.Connected(){}

	var ak ImportedAppKey
	ak.iface = iface
	ak.key = key
	ak.tags = tags

	ak.listenChan = make(chan Result)

	impReq := AurSirAddImportMessage{key, tags}

	iface.out <- impReq

	impRep := <-iface.in

	impMsg, ok := impRep.(AurSirImportAddedMessage)
	if !ok {
		panic("insane runtime!!!")
	}

	ak.importId = impMsg.ImportId
	ak.Connected = impMsg.Exported
	(*iface.imports)[impMsg.ImportId] = &ak
	return &ak
}

//registerResultChan stores a request uuid for ONE2.. or appkey.functioname@tag for MANY2... calls together with a channel so when Results comes in later, it can be routed to the
// appropriate channel.
func (iface *AurSirInterface) registerResultChan(resUuid string, rc chan Result) {
	(*iface.resultChans)[resUuid] = rc
}

//connect initializes the connection to the runtime by sending a DOCK message
func (iface *AurSirInterface) connect() {
	iface.out <- AurSirDockMessage{iface.AppName,[]string{"MSGPACK","JSON"}}
}

func (iface *AurSirInterface) backend() {

	go iface.listener()

	go iface.sender()

}

func (iface *AurSirInterface) sender() {

	log.Println("Opening ougoing ZeroMQ Socket")

	skt, err := zmq4.NewSocket(zmq4.DEALER)

	if err != nil {
		panic("Could not open ZeroMQ socket")
	}

	defer skt.Close()

	skt.SetIdentity(iface.UUID)

	skt.Connect("tcp://localhost:5555")

	log.Println("Outgoing ZeroMQ Socket open")

	for msg := range iface.out {

		var appmsg AppMessage

		appmsg.Encode(msg, "JSON")

		if appmsg.MsgType == DOCK {
			skt.SendMessage([]string{strconv.FormatInt(appmsg.MsgType, 10), appmsg.MsgCodec, string(appmsg.Msg), strconv.FormatInt(iface.port, 10)}, 0)
		} else {
			skt.SendMessage([]string{strconv.FormatInt(appmsg.MsgType, 10), appmsg.MsgCodec, string(appmsg.Msg)}, 0)
		}
	}
	log.Println("Sending LEAVE")
	var lmsg AppMessage

	lmsg.Encode(AurSirLeaveMessage{}, "JSON")

	skt.SendMessage([]string{strconv.FormatInt(lmsg.MsgType, 10), lmsg.MsgCodec, string(lmsg.Msg)}, 0)
	*iface.quit = true
}

func (iface *AurSirInterface) listener() {

	log.Println("Opening Incoming ZeroMQ Socket")

	skt, err := zmq4.NewSocket(zmq4.ROUTER)

	if err != nil {
		panic("Could not open ZeroMQ socket")
	}

	defer skt.Close()

	skt.Bind("tcp://*:" + strconv.FormatInt(iface.port, 10))

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

func (iface *AurSirInterface) processMsg(message []string) {

	msgType, err := strconv.ParseInt(message[1], 10, 64)

	if err != nil {
		return
	}

	switch msgType {

	case DOCKED:
		go pingUdp(iface.UUID)
		iface.connected <- true


	case IMPORT_UPDATED:
		encmsg := []byte(message[3])
		msg := AppMessage{msgType, message[2], encmsg}
		asmsg, err := msg.Decode()
		if err == nil {
			iumsg, ok := asmsg.(AurSirImportUpdatedMessage)
			if ok {
				log.Println(*iface.imports)
				(*iface.imports)[iumsg.ImportId].Connected = iumsg.Exported
			}

		}

	case REQUEST:
		encmsg := []byte(message[3])
		msg := AppMessage{msgType, message[2], encmsg}
		asmsg, err := msg.Decode()

		if err == nil {
			reqmsg, ok := asmsg.(AurSirRequest)

			if ok {
				//this could and should be done more elgently
				<-iface.exportsSem
				for _, exp := range *iface.exports {
					log.Println(exp.key.ApplicationKeyName)
					if exp.key.ApplicationKeyName == reqmsg.AppKeyName {

						exp.Request <- reqmsg

					}
				}
				iface.exportsSem <- struct{}{}
			}

		}
	case RESULT:
		encmsg := []byte(message[3])
		msg := AppMessage{msgType, message[2], encmsg}
		asmsg, err := msg.Decode()
		if err == nil {
			resmsg, ok := asmsg.(AurSirResult)

			if ok {

				var resId string

				if resmsg.CallType == ONE2MANY || resmsg.CallType == ONE2ONE {
					resId = resmsg.Uuid
				} else {
					resId = resmsg.AppKeyName + "." + resmsg.FunctionName
				}

				rc, f := (*iface.resultChans)[resId]

				if f {
					rc <- Result{
						resmsg.Result,
						resmsg.FunctionName,
						resmsg.Uuid,
						resmsg.CallType,
						resmsg.Stream,
						resmsg.StreamFinished,
						resmsg.Codec}

				}

			}

		}

	case CALLCHAIN_ADDED:
		encmsg := []byte(message[3])
		msg := AppMessage{msgType, message[2], encmsg}
		ccmsg, err := msg.Decode()
		if err == nil {
			iface.in <- ccmsg
		}

	default:
		encmsg := []byte(message[3])
		msg := AppMessage{msgType, message[2], encmsg}
		asmsg, err := msg.Decode()

		if err == nil {
			iface.in <- asmsg
		}
	}
}

//Connected returns true if interface successfully connected to a runtime
func (iface *AurSirInterface) Connected() (connected bool) {

	connected = <-iface.connected
	iface.connected <- connected
	return
}

func getRandomPort() int64 {
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
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
