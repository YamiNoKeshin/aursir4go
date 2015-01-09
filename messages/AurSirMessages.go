package messages

import (
	"time"
	"github.com/joernweissenborn/aursir4go/appkey"
	"github.com/joernweissenborn/aursir4go/util"
	"errors"
)

//AurSirMessage represents a generic message
type AurSirMessage interface {
}

const (
	DOCK            = iota
	DOCKED
	LEAVE
	REQUEST
	RESULT
	ADD_EXPORT
	UPDATE_EXPORT
	EXPORT_ADDED
	ADD_IMPORT
	UPDATE_IMPORT
	IMPORT_ADDED
	IMPORT_UPDATED
	LISTEN
	STOP_LISTEN
)

//An DockMessage indicates that an app wants to registered at the engine
type DockMessage struct {
	AppName string
	Codecs []string
}

//An DockMessage indicates that an app is sucessfully registered
type DockedMessage struct {
	Ok bool
}

//An LeaveMessage contains no information since it only serves to indicate that an application is leaving
type LeaveMessage struct{}

//An AddImportMessage contains the respective AppKey struct and a slice with tags
type AddImportMessage struct {
	AppKey appkey.AppKey
	Tags   []string
}

//An UpdateImportMessage contains the import id and the new tag set
type UpdateImportMessage struct {
	ImportId string
	Tags     []string
}

//An ImportAddedMessage contains the import id and a bool flag to indicate if an exporter is connected
type ImportAddedMessage struct {
	ImportId   string
	Exported   bool
	AppKeyName string
	Tags       []string
}

//An ImportAddedMessage contains the import id and a bool flag to indicate if an exporter is connected
type ImportUpdatedMessage struct {
	ImportId string
	Exported bool
}

//An AddExportMessage contains the respective AppKey struct and a slice with tags
type AddExportMessage struct {
	AppKey appkey.AppKey
	Tags   []string
}

//An ExportAddedMessage contains the export id
type ExportAddedMessage struct {
	ExportId string
}

//An UpdateExportMessage contains the export id and the new tag set
type UpdateExportMessage struct {
	ExportId string
	Tags     []string
}

//An ListenMessage contains the import id and the function name.
type ListenMessage struct {
	ImportId     string
	FunctionName string
}

//An StopListenMessage contains the import id and the function name.
type StopListenMessage struct {
	ImportId     string
	FunctionName string
}

//An Request contains the AppKey's name together with the respective Function name. The
type Request struct {
	AppKeyName   string
	FunctionName string
	CallType     int64
	Tags         []string
	Uuid         string
	ImportId     string
	ExportId     string
	Timestamp	time.Time
	Codec        string
	Request      []byte
	Stream	bool
	StreamFinished bool
}



func (r Request) Decode(target interface{}) error {

	codec := util.GetCodec(r.Codec)
	if codec == nil {
		return errors.New("Unknown codec "+r.Codec)
	}

	return codec.Decode(r.Request, &target)
}


type Result struct {
	AppKeyName   string
	FunctionName string
	CallType     int64
	Tags         []string
	Uuid         string
	ImportId     string
	ExportId     string
	Timestamp	time.Time
	Codec        string
	Result       []byte
	Stream	bool
	StreamFinished bool
}


func (r Result) Decode(target interface{}) error {
	codec := util.GetCodec(r.Codec)
	if codec == nil {
		return errors.New("Unknown codec "+r.Codec)
	}
	return codec.Decode(r.Result, &target)
}


