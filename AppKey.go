package aursir4go

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

//func (ak AppKey) Hash() (string){
//hasher := sha1.New()
//hasher.Write([]byte(ak))
//sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
//return sha
//}

type Request struct {
	req      []byte
	Function string
	Uuid     string
	CallType int64
	ImportId string
}

func (r Request) Decode(target interface{}) {
	codec := GetCodec("JSON")
	codec.Decode(&r.req, &target)
}

type Result struct {
	res      []byte
	Function string
	Uuid     string
	CallType int64
}

func (r Result) Decode(target interface{}) {
	codec := GetCodec("JSON")
	codec.Decode(&r.res, &target)
}
