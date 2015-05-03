package aurbaeth

import (
	"github.com/joernweissenborn/stream2go"
)

type AurBaethExport interface {

	Id() string

	Remove()
	UpdateTags([]string)

	Request() AurBaethRequest
	RequestStream()  stream2go.Stream
	RequestChan(function string) (request chan AurBaethRequest)

	Emit(function string, parameter interface {})
}


