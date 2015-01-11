package keys

import "github.com/joernweissenborn/aursir4go/appkey"

var CountStringKey = appkey.AppKey{
	"org.aursir.countstring",
	[]appkey.Function{
		appkey.Function{
			"CountString",
			[]appkey.Data{
				appkey.Data{
					"String",
					appkey.STRING}},
			[]appkey.Data{
				appkey.Data{
					"Size",
					appkey.INT}}}}}

type CountStringReq struct {
	String string
}

type CountStringRes struct {
	Size int64
}
