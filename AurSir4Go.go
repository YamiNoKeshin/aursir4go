package aursir4go

import (
	uuid "github.com/nu7hatch/gouuid"
	"time"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/internal/dock"
	"github.com/joernweissenborn/aursir4go/internal/dock/dockzmq"
	"github.com/joernweissenborn/aursir4go/util"
	"github.com/joernweissenborn/aursir4go/internal"
	"github.com/joernweissenborn/aursir4go/appkey"
)


//The AurSirInterface handles the runtime time connection. It holds methods to create im- and exports.
type AurSirInterface struct {
	AppName string //Name of the application

	UUID string //the applications uuid

	incomingprocessor *internal.Incomingprocessor

	exports map[string]*ExportedAppKey


	imports map[string]*ImportedAppKey



	incoming dock.Incoming
	outgoing dock.Outgoing

	codec string
}

//NewInterface creates and launches a new AurSirInterface and returns a pointer to it.
func NewInterface(AppName string) (iface AurSirInterface, err error) {
	//set the app name
	iface.AppName = AppName

	//generate a UUID
	iface.UUID = generateUuid()

	iface.codec = "JSON"


	iface.incomingprocessor = internal.InitIncomingProcessor()
	iface.exports = map[string]*ExportedAppKey{}

	iface.imports = map[string]*ImportedAppKey{}

	incoming := new(dockzmq.IncomingZmq)
	outgoing, err := incoming.Activate(iface.incomingprocessor)

	if err != nil {
		return
	}

	err = outgoing.Activate(iface.UUID)

	if err != nil {
		return
	}

	iface.incoming = incoming
	iface.outgoing = outgoing

	iface.dock()

	return

}

//Close shuts down the AurSir interface
func (iface *AurSirInterface) Close() {
//	log.Println("Closing out channel")
//	log.Println("out channel closed")
	                        iface.leave()
}
//AddExport adds the specified ApplicationKey and tags and registeres it at the runtime. It returns a pointer to an
// ExportedAppKey which can be userd to handle incoming requests
func (iface *AurSirInterface) AddExport(key appkey.AppKey, tags []string) *ExportedAppKey {

	iface.WaitUntilDocked()
	var ak ExportedAppKey
	ak.iface = iface
	ak.key = key
	ak.tags = tags

	expReq := messages.AddExportMessage{key, tags,""}

	msg, _ := util.GetCodec(iface.codec).Encode(expReq)
	iface.outgoing.Send(messages.ADD_EXPORT,iface.codec,msg)


	ak.exportId = <-iface.incomingprocessor.AddExport

	ak.Request = make(chan messages.Request,1)
	iface.incomingprocessor.RegisterRequestChan(ak.exportId,ak.Request)
	iface.exports[ak.exportId] = &ak

	return &ak
}

//AddImport adds the specified ApplicationKey and tags and registeres it at the runtime. It returns a pointer to an
// ImportedAppKey which can be used to request function calls or call chains and listening to functions
func (iface *AurSirInterface) AddImport(key appkey.AppKey, tags []string) *ImportedAppKey {

	iface.WaitUntilDocked()

	var ak ImportedAppKey
	ak.iface = iface
	ak.key = key
	ak.tags = tags

	ak.listenChan = make(chan messages.Result)
	impReq := messages.AddImportMessage{key, tags,""}
	msg, _ := util.GetCodec(iface.codec).Encode(impReq)
	iface.outgoing.Send(messages.ADD_IMPORT,iface.codec,msg)

	ak.importId = <-iface.incomingprocessor.AddImport

	iface.incomingprocessor.RegisterResultChan(ak.importId,ak.listenChan)
	ak.connected = iface.incomingprocessor.ExportedChans[ak.importId]
	iface.imports[ak.importId] = &ak
	return &ak
}


//dock initializes the connection to the runtime by sending a DOCK message
func (iface AurSirInterface) dock() {
	msg, _ := util.GetCodec(iface.codec).Encode(iface.getdockmsg())
	iface.outgoing.Send(messages.DOCK, iface.codec, msg)
}
//dock initializes the connection to the runtime by sending a DOCK message
func (iface AurSirInterface) leave() {
	iface.outgoing.Send(messages.LEAVE,iface.codec,[]byte{})
}

func (iface AurSirInterface) Docked() bool {
	docked := <- iface.incomingprocessor.Docked
	iface.incomingprocessor.Docked <- docked
	return docked
}
func (iface AurSirInterface) WaitUntilDocked()  {
	docked := false
	for !docked {
		docked = iface.Docked()
		time.Sleep(10 * time.Millisecond)
	}
	return
}
func (iface AurSirInterface) getdockmsg() messages.DockMessage {
	return messages.DockMessage{iface.AppName,[]string{"MSGPACK","JSON"},false}
}





func generateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		//log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}
