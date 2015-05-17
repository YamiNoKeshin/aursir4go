package aurbaeth

import (
	"github.com/joernweissenborn/aursir4go/appkey"
	"github.com/joernweissenborn/future2go"
)

type AurBaethInterface interface {
	Ready() future2go.Future
	AddImport(appkey.AppKey) AurBaethImport
	AddExport(appkey.AppKey) AurBaethExport
	Close()
}
