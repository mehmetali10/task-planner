package postgres_repository

import (
	"context"
	"fmt"

	"github.com/mehmetali10/task-planner/pkg/log"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/pkg/database/postgres"
	"github.com/mehmetali10/task-planner/internal/pkg/database/postgres/tables"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
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
		"Checking if task already exists externalId=%s, provider=%s",
		req.ExternalID,
		req.Provider,
	)

	// Check if a task with the same ExternalID and Provider already exists
	existingTasks, err := postgres.Read[[]payload.CreateTaskResponse, tables.Task](
		ctx,
		map[string]interface{}{
			"ExternalID": req.ExternalID,
			"Provider":   req.Provider,
		},
		1, // Limit to 1 result
		0, // Offset
	)
	if err != nil {
		p.logger.Error(
			"Failed to check existing task externalId=%v, provider=%v: error=%v",
			req.ExternalID,
			req.Provider,
			err,
		)
		return payload.CreateTaskResponse{}, err
	}

	if len(existingTasks) > 0 {
		p.logger.Warn(
			"Task already exists externalId=%v, provider=%v",
			req.ExternalID,
			req.Provider,
		)
		return payload.CreateTaskResponse{}, fmt.Errorf("task with externalId=%v and provider=%v already exists", req.ExternalID, req.Provider)
	}

	p.logger.Trace(
		"Creating new task externalId=%s, provider=%s",
		req.ExternalID,
		req.Provider,
	)
	resp, err := postgres.Create[payload.CreateTaskResponse, tables.Task](ctx, req)
	if err != nil {
		p.logger.Error(
			"Failed to create task externalId=%v, provider=%v: error=%v",
			req.ExternalID,
			req.Provider,
			err,
		)
		return resp, err
	}
	p.logger.Trace(
		"Task created successfully externalId=%v, provider=%v",
		req.ExternalID,
		req.Provider,
	)
	return resp, err
}

// ListTasks implements repository.Repository.
func (p *PostgresRepo) ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error) {
	p.logger.Trace("Listing tasks")
	tasks, err := postgres.Read[[]payload.Task, tables.Task](ctx, map[string]interface{}{}, req.Limit, req.Offset)
	if err != nil {
		p.logger.Error("Failed to list tasks: error=%v", err)
	}
	return payload.ListTasksResponse{Tasks: tasks}, err
}

// ListDevelopers implements repository.Repository.
func (p *PostgresRepo) ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error) {
	p.logger.Trace("Listing developers")
	developers, err := postgres.Read[[]payload.Developer, tables.Developer](ctx, map[string]interface{}{}, 10000, 0)
	if err != nil {
		p.logger.Error("Failed to list developers: error=%v", err)
	}
	return payload.ListDevelopersResponse{Developers: developers}, err
}

// ScheduleAssignments implements repository.Repository.
func (p *PostgresRepo) ScheduleAssignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	return payload.ScheduleAssignmentResponse{}, nil
}
