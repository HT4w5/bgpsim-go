package simulation

import (
	"log"

	"github.com/HT4w5/bgpsim-go/internal/pkg/config"
	"github.com/HT4w5/bgpsim-go/internal/pkg/model"
)

type Simulation struct {
	cfg    *config.Config
	nodes  []*model.Node
	logger *log.Logger
}

func New(cfg *config.Config, logger *log.Logger) *Simulation {
	s := &Simulation{
		cfg:    cfg,
		nodes:  make([]*model.Node, 0, len(cfg.Nodes)),
		logger: logger,
	}

	// Create nodes
	for _, v := range cfg.Nodes {
		if n, err := model.NewNode(v); err != nil {
			logger.Panicf("failed to create node %s: %v", v.Name, err)
		} else {
			s.nodes = append(s.nodes, n)
		}
	}

	return s
}
