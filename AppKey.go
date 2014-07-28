package AurSir4Go


type AppKey struct {

	ApplicationKeyName string

	Functions []Function
}

type Function struct {

	Name string

	Input []Data

	Output []Data

}

type Data struct {

	Name string

	Type int

}





type Request struct {
	req []byte
	Function string
	Uuid string
	CallType int64
	ImportId string
}

func (r Request) Decode(target interface {}){
	codec:= getCodec("JSON")
	codec.decode(&r.req,&target)
}

type Result struct {
	res []byte
	Function string
	Uuid string
	CallType int64
}
func (r Result) Decode(target interface {}){
	codec:= getCodec("JSON")
	codec.decode(&r.res,&target)
}
