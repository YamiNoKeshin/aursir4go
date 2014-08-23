package aursir4go

import (
	"fmt"
	"sort"
	"crypto/md5"
	"encoding/base64"
)

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

func HashAppKey(AppKey AppKey) string{
	funs := make([]string,len(AppKey.Functions))
	for j, f := range AppKey.Functions{
		fstring := f.Name
		inputs := make([]string,len(f.Input))
		for i,in := range f.Input {
			inputs[i] = fmt.Sprintf("%s%d",in.Name,in.Type)
		}
		sort.Strings(inputs)
		outputs := make([]string,len(f.Output))
		for i,out := range f.Output {
			outputs[i] = fmt.Sprintf("%s%d",out.Name,out.Type)
		}
		sort.Strings(inputs)
		for _,s := range inputs {
			fstring = fstring+s
		}
		for _,s := range outputs {
			fstring = fstring+s
		}
		funs[j] = fstring
	}
	sort.Strings(funs)
	keystring := AppKey.ApplicationKeyName
	for _,f := range funs{
		keystring = keystring+f
	}
	hasher := md5.New()
	hasher.Write([]byte(keystring))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

type Request struct {
	req      []byte
	Function string
	Uuid     string
	CallType int64
	ImportId string
	codec string
}

func (r Request) Decode(target interface{}) {
	codec := GetCodec(r.codec)
	codec.Decode(&r.req, &target)
}

type Result struct {
	res      []byte
	Function string
	Uuid     string
	CallType int64
	Stream	bool
	Finished	bool
	codec string
}

func (r Result) Decode(target interface{}) error {
	codec := GetCodec(r.codec)
	return codec.Decode(&r.res, &target)
}
