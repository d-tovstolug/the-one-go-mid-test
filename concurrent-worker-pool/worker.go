package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	Name                  string
	SleepTime             time.Duration
	ErrProbabilityPercent int
}

func NewRandomTask(cfg *Config, name string) *Task {
	sleepTime := rand.Intn(cfg.MaxWaitTimeSeconds-cfg.MinWaitTimeSeconds+1) + cfg.MinWaitTimeSeconds

	return &Task{
		Name:                  name,
		SleepTime:             time.Second * time.Duration(sleepTime),
		ErrProbabilityPercent: cfg.ErrProbabilityPercent,
	}
}

type TaskResult struct {
	Result     int
	ResultErr  error
	WorkerName string
	Task       *Task
}

type Worker interface {
	Run(ctx context.Context) error
}

type worker struct {
	name    string
	tasksCh <-chan *Task
	resCh   chan<- *TaskResult
	logger  *zap.Logger
}

func NewWorker(name string, tasksCh <-chan *Task, resCh chan<- *TaskResult) Worker {
	return &worker{
		name:    name,
		tasksCh: tasksCh,
		resCh:   resCh,
		logger:  zap.L().With(zap.String("service", name)),
	}
}

func (w *worker) Run(ctx context.Context) error {
	w.logger.Info("Starting worker...")
	for {
		select {
		case t := <-w.tasksCh:
			w.processTask(ctx, t)
		case <-ctx.Done():
			w.logger.Info("worker stopped due to ctx canceled")
			return nil
		}
	}
}

func (w *worker) processTask(ctx context.Context, task *Task) {
	select {
	case <-time.After(task.SleepTime):
		res := &TaskResult{
			Task:       task,
			WorkerName: w.name,
		}
		if rand.Intn(100) < task.ErrProbabilityPercent {
			res.ResultErr = fmt.Errorf("task %s processing error", task.Name)
			w.resCh <- res
			return
		}

		res.Result = rand.Int()
		w.resCh <- res
	case <-ctx.Done():
		w.logger.Info("stop processing task due to ctx canceled", zap.String("task", task.Name))
		return
	}
}

type consumer struct {
	resCh   <-chan *TaskResult
	tasksWg *sync.WaitGroup
	logger  *zap.Logger
}

func NewConsumer(resCh <-chan *TaskResult, tasksWg *sync.WaitGroup) Worker {
	return &consumer{
		resCh:   resCh,
		tasksWg: tasksWg,
		logger:  zap.L().With(zap.String("service", "consumer")),
	}
}

func (c *consumer) Run(ctx context.Context) error {
	for {
		select {
		case res := <-c.resCh:
			logger := c.logger.With(zap.String("task", res.Task.Name),
				zap.String("worker", res.WorkerName),
				zap.String("task time", res.Task.SleepTime.String()))
			c.tasksWg.Done()

			if res.ResultErr != nil {
				logger.Error("error processing task. stop processing",
					zap.Error(res.ResultErr))
				return res.ResultErr
			}
			logger.Info("task processed successfully",
				zap.Int("result", res.Result))
		case <-ctx.Done():
			c.logger.Info("ctx closed. stop processing tasks results")
		}
	}
}
