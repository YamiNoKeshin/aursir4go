package aurbaeth

import (
	"github.com/joernweissenborn/stream2go"
)

type AurBaethExport interface {
	Id() string

	Remove()

	UpdateTags([]string)

	ExportFunction(function string) stream2go.Stream

	Emit(function string, parameter interface{})
}
