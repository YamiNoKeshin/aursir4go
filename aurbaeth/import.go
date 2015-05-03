package aurbaeth

import (
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/future2go"
)

type AurBaethImport interface {
	Id() string
	Remove()
	Call(function string, parameter interface {}) (result future2go.Future)
	CallAll(function string, parameter interface {}) (results stream2go.Stream)
	Trigger(function string, parameter interface {})
	TriggerAll(function string, parameter interface {})
	UpdateTags([]string)
	Listen(function string) (results stream2go.Stream)
	Exported() bool
	ExportedChange() (results stream2go.Stream)

}
