package cmd

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
	"time"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
	postgresRepo "github.com/mehmetali10/task-planner/internal/pkg/repository/postgres"
	"github.com/mehmetali10/task-planner/internal/task/migrate"
	"github.com/mehmetali10/task-planner/pkg/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the task planner application",
	Long:  `Start the task planner application, run migrations, and process tasks from providers.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.NewLogger("cli", "trace")

		// 1. Load configuration
		if err := config.LoadConfig(); err != nil {
			logger.Fatal(err.Error())
		}

		// 2. Run migrations
		logger.Info("Running migrations...")
		migrate.MigrateAndSeed(logger)

		// 3. Define providers
		providers := []string{
			"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one",
			"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-two",
		}

		if len(providers) == 0 {
			logger.Fatal("No providers specified.")
		}

		// 4. Initialize repository
		repo := postgresRepo.NewPostgresRepo()

		// 5. Start Worker Pool
		wp := NewWorkerPool(2, repo) // 2 workers
		ctx, cancel := context.WithCancel(context.Background())
		wp.Start(ctx)

		// 6. Create a channel for OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		var wg sync.WaitGroup

		// 7. Run providers in parallel
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

		// 8. After all providers are processed, stop the workers
		go func() {
			wg.Wait() // Wait for providers to finish
			logger.Info("All providers processed, waiting for workers to finish...")
			wp.Stop() // Stop the workers
			cancel()  // Cancel the context
			logger.Info("Worker pool fully stopped. Exiting application.")
			os.Exit(0) // Exit when all tasks are done
		}()

		// 9. Graceful shutdown when receiving OS signal
		<-sigChan
		logger.Info("Received termination signal, shutting down...")
		cancel()
		wp.Stop()

		logger.Info("Application terminated gracefully.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

// fetchAndProcessTasks fetches tasks from a provider and processes them
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

// mapToTask maps raw task data to a CreateTaskRequest
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

// WorkerPool management
type WorkerPool struct {
	taskQueue chan payload.CreateTaskRequest
	workerNum int
	repo      repository.Repository
	wg        sync.WaitGroup // WaitGroup to track workers
	logger    log.Logger
}

// NewWorkerPool creates a new Worker Pool
func NewWorkerPool(workerNum int, repo repository.Repository) *WorkerPool {
	return &WorkerPool{
		taskQueue: make(chan payload.CreateTaskRequest, 100), // 100 buffer size
		workerNum: workerNum,
		repo:      repo,
		logger:    log.NewLogger("worker-pool", "trace"),
	}
}

// Start workers
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1) // Add a worker
		go wp.worker(ctx, i)
	}
}

// Worker function
func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
	defer wp.wg.Done() // Remove from WaitGroup when worker finishes

	wp.logger.Trace("Worker %d started...", workerID)
	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				wp.logger.Trace("Worker %d: Task queue closed, exiting...", workerID)
				return
			}

			resp, err := wp.repo.CreateTask(ctx, task)
			if err != nil {
				wp.logger.Error("Worker %d: Error creating task: %v", workerID, err)
			} else {
				wp.logger.Debug("Worker %d: Task created successfully: %+v", workerID, resp)
			}

		case <-ctx.Done():
			wp.logger.Info("Worker %d stopping...", workerID)
			time.Sleep(time.Millisecond * 100)
			return
		}
	}
}

// SubmitTask submits a task to the WorkerPool
func (wp *WorkerPool) SubmitTask(task payload.CreateTaskRequest) {
	wp.taskQueue <- task
}

// Stop stops the WorkerPool
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue) // Close the queue, workers will start shutting down
	wp.wg.Wait()        // Wait for all workers to finish
}
