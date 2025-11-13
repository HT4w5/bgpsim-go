package model

const (
	defaultLocalAdminCost = 220
	defaultEbgpAdminCost  = 200
	defaultIbgpAdminCost  = 20
)

type NodeConfig struct {
	Name           string               `json:"name"`
	RouterId       int                  `json:"routerId"`
	EbgpAdminCost  int                  `json:"ebgpAdminCost"`
	LocalAdminCost int                  `json:"localAdminCost"`
	MultipathEbgp  int                  `json:"multipathEbgp"`
	LocalAs        int                  `json:"localAs"`
	StaticRoutes   []*StaticRouteConfig `json:"staticRoutes"`
}

type Node struct {
	name           string
	routerId       int
	mainRib        *MainRib
	links          []*Link
	bgpRib         *BgpRib
	localAs        int
	localAdminCost int
	ebgpAdminCost  int
	ibgpAdminCost  int
}

func NewNode(cfg *NodeConfig) (*Node, error) {
	n := &Node{
		name:           cfg.Name,
		routerId:       cfg.RouterId,
		localAs:        cfg.LocalAs,
		links:          []*Link{},
		localAdminCost: defaultLocalAdminCost,
		ebgpAdminCost:  defaultEbgpAdminCost,
		ibgpAdminCost:  defaultIbgpAdminCost,
	}

	// TODO: verify config
	return n, nil
}
