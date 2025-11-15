package model

type BGPConfig struct {
	RouterId           string           `json:"router_id"`
	NetworksAdvertised []string         `json:"networks_advertised"`
	Peers              []*BGPPeerConfig `json:"peers"`
}

type BGP struct {
	routerId string
}
