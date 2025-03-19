package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	postgresRepo "github.com/mehmetali10/task-planner/internal/pkg/repository/postgres"
	"github.com/spf13/cobra"
)

func main() {
	var providers []string

	rootCmd := &cobra.Command{
		Use:   "task-cli",
		Short: "CLI for managing tasks",
		Run: func(cmd *cobra.Command, args []string) {
			if len(providers) == 0 {
				log.Fatal("No providers specified. Use the --providers flag to specify provider URLs.")
			}

			repo := postgresRepo.NewPostgresRepo()

			for _, provider := range providers {
				tasks, err := fetchTasksFromProvider(provider)
				if err != nil {
					log.Printf("Error fetching tasks from provider %s: %v", provider, err)
					continue
				}

				for _, task := range tasks {
					_, err := repo.CreateTask(context.Background(), task)
					if err != nil {
						log.Printf("Error creating task in database: %v", err)
					} else {
						fmt.Printf("Task created successfully: %+v\n", task)
					}
				}
			}
		},
	}

	rootCmd.Flags().StringSliceVar(&providers, "providers", []string{}, "List of provider URLs to fetch tasks from")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func fetchTasksFromProvider(url string) ([]payload.CreateTaskRequest, error) {
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
			log.Printf("Error mapping task: %v", err)
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
	// Handle provider1 format
	if zorluk, ok := raw["zorluk"]; ok {
		return payload.CreateTaskRequest{
			ExternalID: uint(raw["id"].(float64)),
			Name:       fmt.Sprintf("Task %v", raw["id"]),
			Duration:   int(raw["sure"].(float64)),
			Difficulty: int(zorluk.(float64)),
			Provider:   "provider1",
		}, nil
	}

	// Handle provider2 format
	if value, ok := raw["value"]; ok {
		return payload.CreateTaskRequest{
			ExternalID: uint(raw["id"].(float64)),
			Name:       fmt.Sprintf("Task %v", raw["id"]),
			Duration:   int(raw["estimated_duration"].(float64)),
			Difficulty: int(value.(float64)),
			Provider:   "provider2",
		}, nil
	}
	return payload.CreateTaskRequest{}, nil
}
