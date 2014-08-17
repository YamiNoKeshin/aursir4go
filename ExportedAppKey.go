package aursir4go

type ExportedAppKey struct {
	iface    *AurSirInterface
	key      AppKey
	tags     []string
	exportId string
	Request  chan AurSirRequest
	persistenceStrategies map[string] string
}

func (eak ExportedAppKey) Tags() []string {
	return eak.tags
}

func (eak ExportedAppKey) Reply(req *AurSirRequest, res interface{}) error {
	var aursirResult AurSirResult
	aursirResult.AppKeyName = eak.key.ApplicationKeyName
	aursirResult.Codec = "JSON"
	aursirResult.CallType = req.CallType
	aursirResult.Uuid = req.Uuid
	aursirResult.ImportId = req.ImportId
	aursirResult.ExportId = eak.exportId
	aursirResult.Tags = eak.tags
	aursirResult.FunctionName = req.FunctionName

	if strategy,f := eak.persistenceStrategies[req.FunctionName]; f {
		aursirResult.Persistent = true
		aursirResult.PersistenceStrategy = strategy
	}

	codec := GetCodec("JSON")
	result, err := codec.Encode(res)

	if err == nil {
		aursirResult.Result = *result
		eak.iface.out <- aursirResult
	}
	return err
}
func (eak *ExportedAppKey) UpdateTags(NewTags []string) {
	eak.tags = NewTags
	eak.iface.out <- AurSirUpdateExportMessage{eak.exportId, eak.tags}
}

//SetLogging sets the persitence strategy for function calls of the specified function to "log". This overrides all previous
// persitence strategies!
func (eak *ExportedAppKey) SetLogging(FunctionName string){
	eak.persistenceStrategies[FunctionName] = "log"
}
