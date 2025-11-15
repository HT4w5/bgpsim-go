package main

import (
	"log"
	"os"

	"github.com/HT4w5/bgpsim-go/internal/pkg/config"
)

const (
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
)

func main() {
	cfg := config.New()
	cfg.Load("config.json")

	if len(cfg.Log.Output) != 0 {
		if logFile, err := os.OpenFile(cfg.Log.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			log.Printf("failed to open log file %s: %v\n", cfg.Log.Output, err)
			os.Exit(1)
		} else {
			defer logFile.Close()
			log.SetOutput(logFile)
		}
	}
	log.SetFlags(logFlags)
}
