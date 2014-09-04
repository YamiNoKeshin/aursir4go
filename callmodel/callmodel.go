package callmodel

type AurSirCallModel struct {
	RootFork AurSirFork
	Forks []AurSirFork
}

type AurSirFork struct {
	ForkId string
	CallChain []AurSirCall
}

type AurSirCall struct {
	ApplicationKeyName string
	Tags []string
	CallID string
	CallType int64
	ParameterMaps []ParameterMap
	Parameter map[string]interface {}
	ParameterCodec string
	Forks []string
}

type ParameterMap struct{
	OriginId string
	Map map[string] string
}
