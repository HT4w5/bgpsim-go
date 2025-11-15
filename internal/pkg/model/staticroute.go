package model

type StaticRouteConfig struct {
	Network string         `json:"network"`
	NextHop *NextHopConfig `json:"next_hop"`
}
