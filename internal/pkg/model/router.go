package model

type RouterConfig struct {
	Tag     string        `json:"tag"`
	LocalAs uint32        `json:"local_as"`
	BGP     *BGPConfig    `json:"bgp"`
	Static  *StaticConfig `json:"static"`
}

type Router struct {
	tag     string
	localAs uint32
	mainRIB *MainRIB
	bgp     *BGP
}

func NewRouter(cfg *RouterConfig) (*Router, error) {
	r := &Router{
		tag:     cfg.Tag,
		localAs: cfg.LocalAs,
	}
	return r, nil
}
