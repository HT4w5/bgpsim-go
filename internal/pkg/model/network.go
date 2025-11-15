package model

type NetworkConfig struct {
	Name    string          `json:"network_name"`
	Routers []*RouterConfig `json:"routers"`
}
