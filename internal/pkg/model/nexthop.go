package model

import "net/netip"

type NextHopConfig struct {
	IP    string `json:"ip"`
	Iface string `json:"interface"`
}

type NextHop struct {
	IP    netip.Addr
	Iface string
}
