package dock

import "github.com/joernweissenborn/aursir4go/internal"


type Incoming interface {
	Activate(*internal.Incomingprocessor) (Outgoing ,error)
	Deactivate()
}

