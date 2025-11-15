package model

type RouterConfig struct {
	Tag     string        `json:"tag"`
	LocalAs int           `json:"local_as"`
	BGP     *BGPConfig    `json:"bgp"`
	Static  *StaticConfig `json:"static"`
}

type Router struct {
	tag     string
	localAs int
}
