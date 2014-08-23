package aursir4go

var HelloAurSirAppKey = AppKey{
	"org.aursir.helloaursir",
	[]Function{
		Function{
			"SayHello",
			[]Data{
				Data{
					"Greeting",
					STRING}},
			[]Data{
				Data{
					"Answer",
					STRING}}}}}

type SayHelloReq struct {
	Greeting string
}

type SayHelloRes struct {
	Answer string
}
