package repository

import (
	"context"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
)

type Repository interface {
	CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error)
	GetTask(ctx context.Context, req payload.GetTaskRequest) (payload.GetTaskResponse, error)
	ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error)

	ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error)
	ListAssignments(ctx context.Context, req payload.ListAssignmentsRequest) (payload.ListAssignmentsResponse, error)

	ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error)
}
