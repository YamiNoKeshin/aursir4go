package aursir4go

var CountStringKey = AppKey{
	"org.aursir.countstring",
	[]Function{
		Function{
			"CountString",
			[]Data{
				Data{
					"String",
					STRING}},
			[]Data{
				Data{
					"Size",
					INT}}}}}

type CountStringReq struct {
	String string
}

type CountStringRes struct {
	Size int64
}
