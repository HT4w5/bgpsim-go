package simulator

import (
	"fmt"

	"github.com/HT4w5/bgpsim-go/internal/pkg/model"
)

var logPrefix = func(name string) string {
	return fmt.Sprintf("[%s]", name)
}

type Simulator struct {
	cfg     *model.NetworkConfig
	routers []*model.Node
}

func New(cfg *model.NetworkConfig) (*Simulator, error) {
	s := &Simulator{
		cfg:     cfg,
		routers: make([]*model.Node, 0, len(cfg.Routers)),
	}

	// Create nodes
	for _, v := range cfg.Routers {
		if n, err := model.NewNode(v); err != nil {
			return nil, fmt.Errorf("failed to create node %s: %w", v.Name, err)
		} else {
			s.nodes = append(s.nodes, n)
		}
	}

	return s, nil
}

func (s *Simulator) Close() error {
	if err := s.logFile.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Simulator) tickSerial() error {
	return nil
}

func (s *Simulator) tickParallel() error {
	return nil
}
