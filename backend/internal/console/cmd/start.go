package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mehmetali10/task-planner/internal/console/input"
	pvd "github.com/mehmetali10/task-planner/internal/console/provider"
	"github.com/mehmetali10/task-planner/internal/console/worker"
	"github.com/mehmetali10/task-planner/internal/pkg/config"
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

		setupEnvironment()

		if err := config.LoadConfig(); err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Running migrations...")
		migrate.MigrateAndSeed(logger)

		providers := getProviders()
		if len(providers) == 0 {
			logger.Fatal("No providers specified.")
		}

		repo := postgresRepo.NewPostgresRepo()
		wp := worker.NewWorkerPool(len(providers), repo)
		ctx, cancel := context.WithCancel(context.Background())
		wp.Start(ctx)

		processProviders(providers, logger, wp)

		handleShutdown(cancel, wp, logger)
	},
}

func setupEnvironment() {
	envVars := map[string]string{
		"DB_HOST": "localhost",
		// "DB_HOST":     "my_postgres",
		"DB_PORT":     "5432",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "pass",
		"DB_NAME":     "task",
	}

	for key, defaultVal := range envVars {
		value := input.PromptForEnv(key, defaultVal)
		os.Setenv(key, value)
	}
}

func getProviders() []string {
	defaultProviders := []string{
		"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one",
		"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-two",
	}

	var providers []string
	if input.PromptYesNo(fmt.Sprintf("Do you want to use the default providers?\n 1-) %s\n 2-) %s\n (yes/no)", defaultProviders[0], defaultProviders[1])) {
		providers = append(providers, defaultProviders...)
	}

	for {
		newProvider := input.PromptForInput("Enter a new provider URL (or press Enter to continue):")
		if newProvider == "" {
			break
		}
		providers = append(providers, newProvider)
	}

	return providers
}

func processProviders(providers []string, logger log.Logger, wp *worker.WorkerPool) {
	var wg sync.WaitGroup
	wg.Add(len(providers))

	for _, provider := range providers {
		go func(provider string) {
			defer wg.Done()
			if err := pvd.FetchAndProcessTasks(provider, logger, wp); err != nil {
				logger.Error("Error processing tasks from provider %s: %v", provider, err)
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		logger.Info("All providers processed, waiting for workers to finish...")
		wp.Stop()
		logger.Info("Worker pool fully stopped. Exiting application.")
		os.Exit(0)
	}()
}

func handleShutdown(cancel context.CancelFunc, wp *worker.WorkerPool, logger log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Received termination signal, shutting down...")
	cancel()
	wp.Stop()
	logger.Info("Application terminated gracefully.")
}

func init() {
	rootCmd.AddCommand(startCmd)
}
