package aursir4go

import "time"

//AurSirMessage represents a generic message
type AurSirMessage interface {
}

//An AurSirDockMessage indicates that an app wants to registered at the engine
type AurSirDockMessage struct {
	AppName string
	Codecs []string
}

//An AurSirDockMessage indicates that an app is sucessfully registered
type AurSirDockedMessage struct {
}

//An AurSirLeaveMessage contains no information since it only serves to indicate that an application is leaving
type AurSirLeaveMessage struct{}

//An AurSirAddImportMessage contains the respective AppKey struct and a slice with tags
type AurSirAddImportMessage struct {
	AppKey AppKey
	Tags   []string
}

//An AurSirUpdateImportMessage contains the import id and the new tag set
type AurSirUpdateImportMessage struct {
	ImportId string
	Tags     []string
}

//An AurSirImportAddedMessage contains the import id and a bool flag to indicate if an exporter is connected
type AurSirImportAddedMessage struct {
	ImportId   string
	Exported   bool
	AppKeyName string
	Tags       []string
}

//An AurSirImportAddedMessage contains the import id and a bool flag to indicate if an exporter is connected
type AurSirImportUpdatedMessage struct {
	ImportId string
	Exported bool
}

//An AurSirAddExportMessage contains the respective AppKey struct and a slice with tags
type AurSirAddExportMessage struct {
	AppKey AppKey
	Tags   []string
}

//An AurSirExportAddedMessage contains the export id
type AurSirExportAddedMessage struct {
	ExportId string
}

//An AurSirUpdateExportMessage contains the export id and the new tag set
type AurSirUpdateExportMessage struct {
	ExportId string
	Tags     []string
}

//An AurSirListenMessage contains the import id and the function name.
type AurSirListenMessage struct {
	ImportId     string
	FunctionName string
}

//An AurSirStopListenMessage contains the import id and the function name.
type AurSirStopListenMessage struct {
	ImportId     string
	FunctionName string
}

//An AurSirRequest contains the AppKey's name together with the respective Function name. The
type AurSirRequest struct {
	AppKeyName   string
	FunctionName string
	CallType     int64
	Tags         []string
	Uuid         string
	ImportId     string
	Timestamp	time.Time
	Codec        string
	IsFile		 bool
	Persistent	bool
	PersistenceStrategy string //String to determine the strategy (e.g. "log")
	Request      []byte
	Stream	bool
	StreamFinished bool
}

type AurSirResult struct {
	AppKeyName   string
	FunctionName string
	CallType     int64
	Tags         []string
	Uuid         string
	ImportId     string
	ExportId     string
	Timestamp	time.Time
	Codec        string
	IsFile 		bool
	Persistent	bool
	PersistenceStrategy string //String to determine the strategy (e.g. "log")
	Result       []byte
	Stream	bool
	StreamFinished bool
}

//An AurSirCallChain contains the AppKey's name together with the respective Function name. The
type AurSirCallChain struct {
	OriginRequest AurSirRequest
	CallChain     []ChainCall
	FinalImportId string
	FinalCall     ChainCall
}

type AurSirCallChainAddedMessage struct {
	CallChainOk       bool
	InsaneCallIndices []int64
}


