package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	postgresRepo "github.com/mehmetali10/task-planner/internal/pkg/repository/postgres"
	"github.com/mehmetali10/task-planner/internal/task/migrate"
	"github.com/mehmetali10/task-planner/pkg/log"
)

func main() {
	logger := log.NewLogger("cli", "trace")

	if err := config.LoadConfig(); err != nil {
		logger.Fatal(err.Error())
	}

	migrate.MigrateAndSeed(logger)

	providers := []string{
		"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one",
		"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-two",
	}

	if len(providers) == 0 {
		logger.Fatal("No providers specified. Use the --providers flag to specify provider URLs.")
	}

	// Database connection
	repo := postgresRepo.NewPostgresRepo()

	// Start Worker Pool
	wp := NewWorkerPool(2, repo) // 2 workers
	ctx, cancel := context.WithCancel(context.Background())
	wp.Start(ctx)

	// Create a channel for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Run providers in parallel
	wg.Add(len(providers))
	for _, provider := range providers {
		go func(provider string) {
			defer wg.Done()
			err := fetchAndProcessTasks(provider, logger, wp)
			if err != nil {
				logger.Error("Error processing tasks from provider %s: %v", provider, err)
			}
		}(provider)
	}

	// After all providers are processed, stop the workers
	go func() {
		wg.Wait() // Wait for providers to finish
		logger.Info("All providers processed, waiting for workers to finish...")
		wp.Stop() // Stop the workers
		cancel()  // Cancel the context
		logger.Info("Worker pool fully stopped. Exiting application.")
		os.Exit(0) // Exit when all tasks are done
	}()

	// Graceful shutdown when receiving OS signal
	<-sigChan
	logger.Info("Received termination signal, shutting down...")
	cancel()
	wp.Stop()

	logger.Info("Application terminated gracefully.")
}

// Function to fetch data and send it to the WorkerPool
func fetchAndProcessTasks(url string, logger log.Logger, wp *WorkerPool) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response from %s: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var rawTasks []map[string]interface{}
	if err := json.Unmarshal(body, &rawTasks); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, rawTask := range rawTasks {
		task, err := mapToTask(rawTask)
		if err != nil {
			logger.Error("Error mapping task: %v", err)
			continue
		}
		wp.SubmitTask(task)
	}

	return nil
}

func mapToTask(raw map[string]interface{}) (payload.CreateTaskRequest, error) {
	var task payload.CreateTaskRequest

	if zorluk, ok := raw["zorluk"]; ok {
		id, idOk := raw["id"].(float64)
		sure, sureOk := raw["sure"].(float64)
		difficulty, diffOk := zorluk.(float64)

		if !idOk || !sureOk || !diffOk {
			return task, errors.New("type assertion failed for provider1 fields")
		}

		return payload.CreateTaskRequest{
			ExternalID: uint(id),
			Name:       fmt.Sprintf("Task %v", uint(id)),
			Duration:   int(sure),
			Difficulty: int(difficulty),
			Provider:   "provider1",
		}, nil
	}

	if value, ok := raw["value"]; ok {
		id, idOk := raw["id"].(float64)
		estimatedDuration, durOk := raw["estimated_duration"].(float64)
		difficulty, diffOk := value.(float64)

		if !idOk || !durOk || !diffOk {
			return task, errors.New("type assertion failed for provider2 fields")
		}

		return payload.CreateTaskRequest{
			ExternalID: uint(id),
			Name:       fmt.Sprintf("Task %v", uint(id)),
			Duration:   int(estimatedDuration),
			Difficulty: int(difficulty),
			Provider:   "provider2",
		}, nil
	}

	return task, errors.New("unknown provider format")
}
