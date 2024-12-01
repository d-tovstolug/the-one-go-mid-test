package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
)

func runReadFile(ctx context.Context, filename string, skipHeader bool, resCh chan<- []string, wg *sync.WaitGroup) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s, %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	totalLines := 0

	for {
		line, err := reader.Read()
		switch {
		case errors.Is(err, io.EOF):
			zap.L().Info("file read successfully", zap.Int("total lines", totalLines))
			return nil
		case err != nil:
			return fmt.Errorf("error parsing line %v", err)
		}

		if totalLines == 0 && skipHeader {
			totalLines++
			continue
		}

		wg.Add(1)
		select {
		case resCh <- line:
			totalLines++
		case <-ctx.Done():
			wg.Done()
			return nil
		}
	}
}

func runFileProcessing(ctx context.Context, dataCh <-chan []string, wg *sync.WaitGroup) {
	logger := zap.L().With(zap.String("worker id", uuid.NewString()))
	totalProcessed := 0

	logger.Info("start processing")

	for {
		select {
		case line := <-dataCh:
			processCSVLine(line)
			totalProcessed++
			wg.Done()
		case <-ctx.Done():
			logger.Info("finish processing", zap.Int("total processed", totalProcessed))
			return
		}
	}
}

func processCSVLine(line []string) {
	data, err := parseCSVLine(line)
	if err != nil {
		zap.L().Error("error parsing line", zap.Any("line", line), zap.Error(err))
		return
	}

	zap.L().Debug("read line", zap.Any("data", data))
}
