package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	postgresRepo "github.com/mehmetali10/task-planner/internal/pkg/repository/postgres"
	"github.com/mehmetali10/task-planner/internal/task/migrate"
	"github.com/mehmetali10/task-planner/pkg/log"
)

func main() {
	logger := log.NewLogger("cli", "debug")

	if err := config.LoadConfig(); err != nil {
		logger.Fatal(err.Error())
	}

	migrate.MigrateAndSeed(logger)

	var providers []string = []string{
		"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one",
		"https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-two",
	}

	if len(providers) == 0 {
		logger.Fatal("No providers specified. Use the --providers flag to specify provider URLs.")
	}

	repo := postgresRepo.NewPostgresRepo()

	for _, provider := range providers {
		tasks, err := fetchTasksFromProvider(provider, logger)
		if err != nil {
			logger.Error("Error fetching tasks from provider %s: %v", provider, err)
			continue
		}

		for _, task := range tasks {
			_, err := repo.CreateTask(context.Background(), task)
			if err != nil {
				logger.Error("Error creating task in database: %v", err)
			} else {
				logger.Info("Task created successfully: %+v", task)
			}
		}
	}

}

func fetchTasksFromProvider(url string, logger log.Logger) ([]payload.CreateTaskRequest, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from %s: %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rawTasks []map[string]interface{}
	if err := json.Unmarshal(body, &rawTasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var tasks []payload.CreateTaskRequest
	for _, rawTask := range rawTasks {
		task, err := mapToTask(rawTask)
		if err != nil {
			logger.Error("Error mapping task: %v", err)
			continue
		}
		tasks = append(tasks, payload.CreateTaskRequest{
			ExternalID: task.ExternalID,
			Name:       task.Name,
			Duration:   task.Duration,
			Difficulty: task.Difficulty,
			Provider:   task.Provider,
		})
	}

	return tasks, nil
}

func mapToTask(raw map[string]interface{}) (payload.CreateTaskRequest, error) {
	var task payload.CreateTaskRequest

	// Handle provider1 format
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

	// Handle provider2 format
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
