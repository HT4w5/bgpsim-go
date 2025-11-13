package simulator

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/HT4w5/bgpsim-go/internal/pkg/config"
	"github.com/HT4w5/bgpsim-go/internal/pkg/model"
)

const (
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
)

var logPrefix = func(name string) string {
	return fmt.Sprintf("[%s]", name)
}

type Simulator struct {
	cfg     *config.Config
	nodes   []*model.Node
	logFile *os.File
	logger  *log.Logger
}

func New(cfg *config.Config) (*Simulator, error) {
	s := &Simulator{
		cfg:   cfg,
		nodes: make([]*model.Node, 0, len(cfg.Nodes)),
	}

	// Log to file
	var logWriter io.Writer
	if len(cfg.Log.Output) != 0 {
		if logFile, err := os.OpenFile(cfg.Log.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", cfg.Log.Output, err)
		} else {
			logWriter = logFile
		}
	} else {
		logWriter = os.Stdout
	}
	s.logger = log.New(logWriter, logPrefix(cfg.Name), logFlags)

	// Create nodes
	for _, v := range cfg.Nodes {
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
