package model

import "fmt"

type BGPConfig struct {
	RouterId           string           `json:"router_id"`
	NetworksAdvertised []string         `json:"networks_advertised"`
	Peers              []*BGPPeerConfig `json:"peers"`
}

type BGP struct {
	routerId string
	bgpRIB   *BGP
	peers    []*BGPPeer
}

func NewBGP(cfg *BGPConfig) (*BGP, error) {
	b := &BGP{
		routerId: cfg.RouterId,
	}

	// Create peers
	for _, v := range cfg.Peers {
		p, err := NewBGPPeer(v)
		if err != nil {
			return nil, fmt.Errorf("failed to create peer %s: %w", v.Tag, err)
		}
		b.peers = append(b.peers, p)
	}

	// Create BGPRIB

	return b, nil
}
