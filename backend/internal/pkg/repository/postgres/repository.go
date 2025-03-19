package postgres_repository

import (
	"context"

	"github.com/mehmetali10/task-planner/pkg/log"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
	"github.com/mehmetali10/task-planner/internal/task/config"
	"github.com/mehmetali10/task-planner/pkg/database/postgres"
	"github.com/mehmetali10/task-planner/pkg/database/postgres/tables"
)

type PostgresRepo struct {
	logger log.Logger
}

func NewPostgresRepo() repository.Repository {
	logger := log.NewLogger("postgres-repository", config.GetApp().RepositoryLogLevel)
	logger.Info("Repository instance creating")
	return &PostgresRepo{
		logger: logger,
	}
}

// CreateTask implements repository.Repository.
func (p *PostgresRepo) CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error) {
	p.logger.Trace(
		"Creating new task externalId=%s, provider=%s",
		req.ExternalID,
		req.Provider,
	)
	resp, err := postgres.Create[payload.CreateTaskResponse, tables.Task](ctx, req)
	if err != nil {
		p.logger.Error(
			"Failed to create task externalId=%s, provider=%s: error=%v",
			req.ExternalID,
			req.Provider,
			err,
		)
		return resp, err
	}
	p.logger.Trace(
		"Task created successfully externalId=%s, provider=%s",
		req.ExternalID,
		req.Provider,
	)
	return resp, err
}

// ListTasks implements repository.Repository.
func (p *PostgresRepo) ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error) {
	p.logger.Trace("Listing tasks")
	resp, err := postgres.Read[payload.ListTasksResponse, tables.Task](ctx, map[string]interface{}{}, req.Limit, req.Offset)
	if err != nil {
		p.logger.Error("Failed to list tasks: error=%v", err)
	}
	return resp, err
}

// ListDevelopers implements repository.Repository.
func (p *PostgresRepo) ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error) {
	p.logger.Trace("Listing developers")
	resp, err := postgres.Read[payload.ListDevelopersResponse, tables.Developer](ctx, map[string]interface{}{}, 0, 1000)
	if err != nil {
		p.logger.Error("Failed to list developers: error=%v", err)
	}
	return resp, err
}

// ScheduleAssaignments implements repository.Repository.
func (p *PostgresRepo) ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	return payload.ScheduleAssignmentResponse{}, nil
}
