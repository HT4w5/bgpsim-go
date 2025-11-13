package model

type BgpRib struct {
}

type BgpTopologyConfig struct {
	Nodes []*BgpNodeConfig `json:"nodes"`
	Edges []*BgpEdgeConfig `json:"edges"`
}

type BgpNodeConfig struct {
	Hostname string `json:"hostname"`
	Prefix   string `json:"prefix"`
}

type BgpEdgeConfig struct {
	Source *BgpEdgeNodeConfig `json:"source"`
	Target *BgpEdgeNodeConfig `json:"target"`
}

type BgpEdgeNodeConfig struct {
	Hostname  string `json:"hostname"`
	Prefix    string `json:"prefix"`
	Interface string `json:"interface"`
}
