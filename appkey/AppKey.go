package appkey

import (

	yaml "gopkg.in/yaml.v2"

	"github.com/joernweissenborn/aursir4go/util"
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
	if yaml.Unmarshal([]byte(YAML),&AppKey) !=nil {
		panic("Insane Appkey")
	}
	return
}

