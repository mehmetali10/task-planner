package postgres_repository

import (
	"context"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
	"github.com/mehmetali10/task-planner/pkg/database/postgres"
	"github.com/mehmetali10/task-planner/pkg/database/postgres/tables"
)

type PostgresRepo struct{}

func NewPostgresRepo() repository.Repository {
	return &PostgresRepo{}
}

// CreateTask implements repository.Repository.
func (p *PostgresRepo) CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error) {
	resp, err := postgres.Create[payload.CreateTaskResponse, tables.Task](ctx, req)
	return resp, err
}

// GetTask implements repository.Repository.
func (p *PostgresRepo) GetTask(ctx context.Context, req payload.GetTaskRequest) (payload.GetTaskResponse, error) {
	resp, err := postgres.Read[payload.GetTaskResponse, tables.Task](ctx, map[string]interface{}{"id": req.ID})
	return resp, err
}

// ListTasks implements repository.Repository.
func (p *PostgresRepo) ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error) {
	resp, err := postgres.Read[payload.ListTasksResponse, tables.Task](ctx, map[string]interface{}{})
	return resp, err
}

// ScheduleAssaignments implements repository.Repository.
func (p *PostgresRepo) ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	panic("unimplemented")
}

// ListAssignments implements repository.Repository.
func (p *PostgresRepo) ListAssignments(ctx context.Context, req payload.ListAssignmentsRequest) (payload.ListAssignmentsResponse, error) {
	resp, err := postgres.Read[payload.ListAssignmentsResponse, tables.Assignment](ctx, map[string]interface{}{})
	return resp, err
}

// ListDevelopers implements repository.Repository.
func (p *PostgresRepo) ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error) {
	resp, err := postgres.Read[payload.ListDevelopersResponse, tables.Developer](ctx, nil)
	return resp, err
}
