package appkey

import (

	yaml "gopkg.in/yaml.v2"

	"github.com/joernweissenborn/aursir4go/util"
	"fmt"
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


func (AppKey *AppKey) CreateFromJson(JSON string) error {
	codec := util.GetCodec("JSON")
	return codec.Decode([]byte(JSON), AppKey)
}

func (AppKey *AppKey) CreateFromYaml(YAML string)  {
	if err := yaml.Unmarshal([]byte(YAML),&AppKey); err !=nil {
		panic(fmt.Sprint("Insane Appkey",err))
	}
	return
}



func AppKeyFromJson(JSON string) (appkey AppKey) {
	codec := util.GetCodec("JSON")
	err := codec.Decode([]byte(JSON), &appkey)
	if err!=nil {
		panic(fmt.Sprint("Insane Appkey",err))
	}
	return 
}

func AppKeyFromYaml(YAML string)  (appkey AppKey) {
	if err := yaml.Unmarshal([]byte(YAML),&appkey); err !=nil {
		panic(fmt.Sprint("Insane Appkey",err))
	}
	return
}

