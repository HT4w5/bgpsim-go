package model

import "net/netip"

type BGPOrigin int

const (
	BGPOriginIGP BGPOrigin = iota
	BGPOriginEGP
	BGPOriginIncomplete
)

type BGPPath struct {
	Prefix        netip.Prefix
	NextHop       NextHop
	AsPath        *AsPath
	Origin        BGPOrigin
	LocalPerf     uint32
	MultiExitDesc uint32
	Weight        uint32
}
