package service

import (
	"context"

	"github.com/mehmetali10/task-planner/pkg/log"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
)

type Service interface {
	CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error)
	ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error)

	ScheduleAssignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error)

	ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
}

func NewService(repository repository.Repository) Service {
	logger := log.NewLogger("service", config.GetApp().ServiceLogLevel)
	logger.Trace("Service instance created")
	return &service{
		repository: repository,
		logger:     logger,
	}
}

// CreateTask implements Service.
func (s *service) CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error) {
	s.logger.Trace(
		"Creating new task externalId=%v, provider=%v",
		req.ExternalID,
		req.Provider,
	)
	resp, err := s.repository.CreateTask(ctx, req)
	if err != nil {
		s.logger.Error(
			"Failed to create task externalId=%v, provider=%v: error=%v",
			req.ExternalID,
			req.Provider,
			err,
		)
		return resp, err
	}
	s.logger.Trace(
		"Task created successfully externalId=%v, provider=%v",
		req.ExternalID,
		req.Provider,
	)
	return resp, nil
}

// ListDevelopers implements Service.
func (s *service) ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error) {
	s.logger.Trace("Listing developers")
	resp, err := s.repository.ListDevelopers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list developers: error=%v", err)
		return resp, err
	}
	s.logger.Trace("Developers listed successfully")
	return resp, nil
}

// ListTasks implements Service.
func (s *service) ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error) {
	s.logger.Trace("Listing tasks")
	resp, err := s.repository.ListTasks(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list tasks: error=%v", err)
		return resp, err
	}
	s.logger.Trace("Tasks listed successfully")
	return resp, nil
}
