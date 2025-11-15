package model

type BGPPeerConfig struct {
	Tag      string `json:"tag"`
	RemoteAs int    `json:"remote_as"`
	LocalIP  string `json:"local_ip"`
	RemoteIP string `json:"remote_ip"`
}
