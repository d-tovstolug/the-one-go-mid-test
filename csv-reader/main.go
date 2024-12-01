package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	logger, err := zap.NewProduction()
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	}
	//defer logger.Sync()
	zap.ReplaceGlobals(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	processFile(ctx, cfg)
}

func processFile(ctx context.Context, cfg *Config) {
	wg := &sync.WaitGroup{}
	dataCh := make(chan []string)

	for range cfg.WorkersNum {
		go runFileProcessing(ctx, dataCh, wg)
	}

	if err := runReadFile(ctx, cfg.FileName, cfg.SkipHeader, dataCh, wg); err != nil {
		zap.L().Fatal("error reading file", zap.Error(err))
	}

	wg.Wait()
}
