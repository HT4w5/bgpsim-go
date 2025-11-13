package model

const (
	defaultLocalAdminCost = 220
	defaultEbgpAdminCost  = 200
	defaultIbgpAdminCost  = 20
)

type NodeConfig struct {
	Name           string                   `json:"name"`
	RouterId       string                   `json:"routerId"`
	EbgpAdminCost  int                      `json:"ebgpAdminCost"`
	LocalAdminCost int                      `json:"localAdminCost"`
	MultipathEbgp  int                      `json:"multipathEbgp"`
	NeighborAses   []int                    `json:"neighborAses"`
	StaticRoutes   []*NodeStaticRouteConfig `json:"staticRoutes"`
}

type Node struct {
	name           string
	id             int
	cfg            *NodeConfig
	mainRib        *MainRib
	links          []*Link
	bgpRib         *BgpRib
	localAs        int
	localAdminCost int
	ebgpAdminCost  int
	ibgpAdminCost  int
}

func NewNode() *Node {
	return &Node{
		links:          []*Link{},
		localAdminCost: defaultLocalAdminCost,
		ebgpAdminCost:  defaultEbgpAdminCost,
		ibgpAdminCost:  defaultIbgpAdminCost,
	}
}

type NodeStaticRouteConfig struct {
	Network            string `json:"network"`
	NextHopIp          string `json:"nextHopIp"`
	NextHopInterface   string `json:"nextHopInterface"`
	AdministrativeCost int    `json:"administrativeCost"`
	Tag                int    `json:"tag"`
	Weight             int    `json:"weight"`
	Metric             int    `json:"metric"`
	Class              string `json:"class"`
}
