package AurSir4Go


//AurSirMessage represents a generic message
type AurSirMessage interface {

}

//An AurSirDockMessage indicates that an app wants to registered at the engine
type AurSirDockMessage struct {
	AppName string
}

//An AurSirDockMessage indicates that an app is sucessfully registered
type AurSirDockedMessage struct {
}

//An AurSirLeaveMessage contains no information since it only serves to indicate that an application is leaving
type AurSirLeaveMessage struct {}

//An AurSirAddImportMessage contains the respective AppKey struct and a slice with tags
type AurSirAddImportMessage struct {
	AppKey AppKey
	Tags []string
}

//An AurSirImportAddedMessage contains the import id and a bool flag to indicate if an exporter is connected
type AurSirImportAddedMessage struct {
	ImportId string
	Exported bool
}
//An AurSirImportAddedMessage contains the import id and a bool flag to indicate if an exporter is connected
type AurSirImportUpdatedMessage struct {
	ImportId string
	Exported bool
}

//An AurSirAddExportMessage contains the respective AppKey struct and a slice with tags
type AurSirAddExportMessage struct {
	AppKey AppKey
	Tags []string
}

//An AurSirExportAddedMessage contains the export id
type AurSirExportAddedMessage struct {
	ExportId string
}

//An AurSirListenMessage contains the import id and the function name.
type AurSirListenMessage struct {
	ImportId string
	FunctionName string
}

//An AurSirStopListenMessage contains the import id and the function name.
type AurSirStopListenMessage struct {
	ImportId string
	FunctionName string
}

//An AurSirRequest contains the AppKey's name together with the respective Function name. The
type AurSirRequest struct {
	AppKeyName string
	FunctionName string
	CallType int64
	Tags []string
	Uuid string
	Codec string
	Request []byte
}

type AurSirResult struct {
	AppKeyName string
	FunctionName string
	CallType int64
	Tags []string
	Uuid string
	Codec string
	Result []byte
}

const(
	ONE2ONE = iota
	MANY2ONE = iota
	ONE2MANY = iota
	MANY2MANY = iota
)
