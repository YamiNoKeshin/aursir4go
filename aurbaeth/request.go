package aurbaeth

type AurBaethRequest interface {
	Decode(interface {})
	Reply(params interface {})

}
