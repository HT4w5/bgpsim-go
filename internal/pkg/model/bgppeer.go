package model

import (
	"fmt"
	"net/netip"
)

type BGPPeerConfig struct {
	Tag      string `json:"tag"`
	RemoteAs uint32 `json:"remote_as"`
	LocalIP  string `json:"local_ip"`
	RemoteIP string `json:"remote_ip"`
}

type BGPPeer struct {
	tag      string
	remoteAs uint32
	localIP  netip.Addr
	remoteIP netip.Addr
}

func NewBGPPeer(cfg *BGPPeerConfig) (*BGPPeer, error) {
	p := &BGPPeer{
		tag:      cfg.Tag,
		remoteAs: cfg.RemoteAs,
	}

	if l, err := netip.ParseAddr(cfg.LocalIP); err != nil {
		return nil, fmt.Errorf("invalid local_ip %s: %w", cfg.LocalIP, err)
	} else {
		p.localIP = l
	}

	if r, err := netip.ParseAddr(cfg.RemoteIP); err != nil {
		return nil, fmt.Errorf("invalid remote_ip %s: %w", cfg.LocalIP, err)
	} else {
		p.remoteIP = r
	}

	return p, nil
}
