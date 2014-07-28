package AurSir4Go

type ExportedAppKey struct{
	iface *AurSirInterface
	key AppKey
	tags []string
	exportId string
	Request chan Request
}

func (eak ExportedAppKey) Tags() []string{
	return eak.tags
}

func (eak ExportedAppKey) Reply(req *Request, res interface {}) error {
	var aursirResult AurSirResult
	aursirResult.AppKeyName = eak.key.ApplicationKeyName
	aursirResult.Codec = "JSON"
	aursirResult.CallType = req.CallType
	aursirResult.Uuid = req.Uuid
	aursirResult.ImportId = req.ImportId
	aursirResult.ExportId = eak.exportId
	aursirResult.Tags = eak.tags
	aursirResult.FunctionName = req.Function
	codec:= getCodec("JSON")
	result, err := codec.encode(res)

	if err ==nil {
		aursirResult.Result = *result
		eak.iface.out <- aursirResult
	}
	return err
}
func (eak *ExportedAppKey) UpdateTags(NewTags []string){
	eak.tags = NewTags
	eak.iface.out <- AurSirUpdateExportMessage{eak.exportId,eak.tags}
}
