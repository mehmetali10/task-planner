package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
		logger := log.NewLogger("cli", "error")

		// 1. Prompt user for environment variables
		envVars := map[string]string{
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASSWORD": "pass",
			"DB_NAME":     "task",
		}

		for key, defaultVal := range envVars {
			value := promptForEnv(key, defaultVal)
			os.Setenv(key, value)
		}

		// 2. Load configuration
		if err := config.LoadConfig(); err != nil {
			logger.Fatal(err.Error())
		}

		// 3. Run migrations
		logger.Info("Running migrations...")
		migrate.MigrateAndSeed(logger)

		// 4. Ask user about default providers
		defaultProviders := []string{
			"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one",
			"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-two",
		}

		useDefaults := promptYesNo(fmt.Sprintf("Do you want to use the default providers?\n 1-) %s\n 2-) %s\n (yes/no)", defaultProviders[0], defaultProviders[1]))
		var providers []string
		if useDefaults {
			providers = append(providers, defaultProviders...)
		}

		// 5. Ask user for additional providers
		for {
			newProvider := promptForInput("Enter a new provider URL (or press Enter to continue):")
			if newProvider == "" {
				break
			}
			providers = append(providers, newProvider)
		}

		if len(providers) == 0 {
			logger.Fatal("No providers specified.")
		}

		// 6. Initialize repository
		repo := postgresRepo.NewPostgresRepo()

		// 7. Start Worker Pool
		wp := NewWorkerPool(len(providers), repo) // Dynamic worker count
		ctx, cancel := context.WithCancel(context.Background())
		wp.Start(ctx)

		// 8. Create a channel for OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		var wg sync.WaitGroup

		// 9. Run providers in parallel
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

		// 10. After all providers are processed, stop the workers
		go func() {
			wg.Wait()
			logger.Info("All providers processed, waiting for workers to finish...")
			wp.Stop()
			cancel()
			logger.Info("Worker pool fully stopped. Exiting application.")
			os.Exit(0)
		}()

		// 11. Graceful shutdown when receiving OS signal
		<-sigChan
		logger.Info("Received termination signal, shutting down...")
		cancel()
		wp.Stop()

		logger.Info("Application terminated gracefully.")
	},
}

func promptForEnv(key, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter value for %s (default: %s): ", key, defaultVal)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultVal
		}

		if isValidEnvVar(key, input) {
			return input
		}

		fmt.Println("Invalid value. Please enter a valid input.")
	}
}

func isValidEnvVar(key, value string) bool {
	switch key {
	case "DB_PORT":
		_, err := strconv.Atoi(value)
		return err == nil
	default:
		return len(value) > 0
	}
}

func promptForInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptYesNo(message string) bool {
	for {
		input := promptForInput(message + " ")
		input = strings.ToLower(input)
		if input == "yes" || input == "y" || input == "" {
			return true
		} else if input == "no" || input == "n" {
			return false
		}
		fmt.Println("Please enter 'yes' or 'no'.")
	}
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
		task, err := mapToTask(rawTask, url)
		if err != nil {
			logger.Error("Error mapping task: %v", err)
			continue
		}
		wp.SubmitTask(task)
	}

	return nil
}

// mapToTask maps raw task data to a CreateTaskRequest
func mapToTask(raw map[string]interface{}, provider string) (payload.CreateTaskRequest, error) {
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
			Provider:   provider,
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
			Provider:   provider,
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
		logger:    log.NewLogger("worker-pool", "error"),
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

	wp.logger.Error("Worker %d started...", workerID)
	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				wp.logger.Error("Worker %d: Task queue closed, exiting...", workerID)
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
	wp.logger.Info("Stopping worker pool... Waiting for remaining tasks.")

	time.Sleep(time.Second * 2) // Wait for 2 seconds for remaining tasks to finish

	close(wp.taskQueue) // Close the queue, workers will start shutting down
	wp.wg.Wait()        // Wait for all workers to finish
}
