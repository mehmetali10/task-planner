package service

import (
	"context"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
)

type Service interface {
	CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error)
	GetTask(ctx context.Context, req payload.GetTaskRequest) (payload.GetTaskResponse, error)
	ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error)

	ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error)
	ListAssignments(ctx context.Context, req payload.ListAssignmentsRequest) (payload.ListAssignmentsResponse, error)

	ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error)
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{
		repository: repository,
	}
}

// ScheduleAssaignments implements Service.
func (s *service) ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	panic("unimplemented")
}

// CreateTask implements Service.
func (s *service) CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error) {
	panic("unimplemented")
}

// GetTask implements Service.
func (s *service) GetTask(ctx context.Context, req payload.GetTaskRequest) (payload.GetTaskResponse, error) {
	panic("unimplemented")
}

// ListAssignments implements Service.
func (s *service) ListAssignments(ctx context.Context, req payload.ListAssignmentsRequest) (payload.ListAssignmentsResponse, error) {
	panic("unimplemented")
}

// ListDevelopers implements Service.
func (s *service) ListDevelopers(ctx context.Context, req payload.ListDevelopersRequest) (payload.ListDevelopersResponse, error) {
	panic("unimplemented")
}

// ListTasks implements Service.
func (s *service) ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error) {
	panic("unimplemented")
}
