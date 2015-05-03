package aurbaeth

import (
	"github.com/joernweissenborn/future2go"
	"github.com/joernweissenborn/aursir4go/appkey"
)

type AurBaethInterface interface {
	Ready() future2go.Future
	AddImport(appkey.AppKey) (AurBaethImport)
    AddExport(appkey.AppKey) AurBaethExport
	Close()

}

