package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal("config parsing error", err)
	}

	logger, err := zap.NewProduction()
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("logger init error", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	sourceTasksCh := make(chan *Task)
	resTasksCh := make(chan *TaskResult)
	tasksWg := &sync.WaitGroup{}
	errGrp, grpCtx := errgroup.WithContext(ctx)

	go runTasks(grpCtx, errGrp, cfg, sourceTasksCh, resTasksCh, tasksWg)
	produceTasks(grpCtx, cfg, sourceTasksCh, tasksWg)

	tasksWg.Wait()
	zap.L().Info("task processing finished")
}

func runTasks(ctx context.Context, errGrp *errgroup.Group, cfg *Config,
	sourceTasksCh chan *Task, resTasksCh chan *TaskResult,
	tasksWg *sync.WaitGroup) {

	for i := range cfg.WorkersNum {
		workerName := fmt.Sprintf("worker-%d", i)
		errGrp.Go(func() error {
			return NewWorker(workerName, sourceTasksCh, resTasksCh).Run(ctx)
		})
	}
	errGrp.Go(func() error {
		return NewConsumer(resTasksCh, tasksWg).Run(ctx)
	})

	if err := errGrp.Wait(); err != nil {
		zap.L().Fatal("processing error", zap.Error(err))
	}
}

func produceTasks(ctx context.Context, cfg *Config, sourceTasksCh chan<- *Task, tasksWg *sync.WaitGroup) {
	for i := range cfg.TasksNum {
		t := NewRandomTask(cfg, fmt.Sprintf("task-%d", i))

		select {
		case sourceTasksCh <- t:
			tasksWg.Add(1)
		case <-ctx.Done():
			zap.L().Info("stop producing tasks due to ctx canceled", zap.Int("total produced", i))
			return
		}
	}
}
