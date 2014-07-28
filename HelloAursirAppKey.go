package AurSir4Go

var Testkey = AppKey{
	"org.aursir.helloaursir",
	[]Function{
		Function{
			"SayHello",
			[]Data{
				Data{
					"Greeting",
					1}},
			[]Data{
				Data{
					"Answer",
					1}}}}}

type SayHelloReq struct {
	Greeting string
}

type SayHelloRes struct {
	Answer string
}
