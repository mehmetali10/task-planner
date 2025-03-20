package repository

import (
	"context"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
)

type Repository interface {
	CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error)
	ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error)
	ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error)
}
