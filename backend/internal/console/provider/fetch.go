package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"

	"github.com/mehmetali10/task-planner/internal/console/worker"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/pkg/log"
)

// fetchAndProcessTasks fetches tasks from a provider and processes them
func FetchAndProcessTasks(url string, logger log.Logger, wp *worker.WorkerPool) error {
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
