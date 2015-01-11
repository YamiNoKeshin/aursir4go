package keys

import "github.com/joernweissenborn/aursir4go/appkey"

var HelloAurSirAppKey = appkey.AppKey{
	"org.aursir.helloaursir",
	[]appkey.Function{
		appkey.Function{
			"SayHello",
			[]appkey.Data{
				appkey.Data{
					"Greeting",
					appkey.STRING}},
			[]appkey.Data{
				appkey.Data{
					"Answer",
					appkey.STRING}}}}}

type SayHelloReq struct {
	Greeting string
}

type SayHelloRes struct {
	Answer string
}
