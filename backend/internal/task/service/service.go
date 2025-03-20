package service

import (
	"context"
	"sort"

	"github.com/mehmetali10/task-planner/pkg/log"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/repository"
)

type Service interface {
	CreateTask(ctx context.Context, req payload.CreateTaskRequest) (payload.CreateTaskResponse, error)
	ListTasks(ctx context.Context, req payload.ListTasksRequest) (payload.ListTasksResponse, error)

	ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error)

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

func (s *service) ScheduleAssaignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	s.logger.Trace("Scheduling assignments")
	tasksResp, err := s.repository.ListTasks(ctx, payload.ListTasksRequest{})
	if err != nil {
		s.logger.Error("Failed to list tasks: error=%v", err)
		return payload.ScheduleAssignmentResponse{}, err
	}

	developersResp, err := s.repository.ListDevelopers(ctx, payload.ListDevelopersRequest{})
	if err != nil {
		s.logger.Error("Failed to list developers: error=%v", err)
		return payload.ScheduleAssignmentResponse{}, err
	}

	tasks := tasksResp.Tasks
	developers := developersResp.Developers

	const weeklyWorkHours = 45
	assignments := []payload.Assignment{}
	totalWeeks := 0

	// Sort tasks by estimated hours (longest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Duration > tasks[j].Duration
	})

	// Assign tasks to developers
	for len(tasks) > 0 {
		assignment := payload.Assignment{
			DeveloperTasks: []payload.DeveloperTaskAssignment{},
		}

		// Reset developer workloads at the beginning of each week
		developerWorkloads := make(map[uint]int)
		for _, dev := range developers {
			developerWorkloads[dev.ID] = 0
		}

		for _, dev := range developers {
			remainingHours := weeklyWorkHours - developerWorkloads[dev.ID]
			if remainingHours <= 0 {
				continue
			}

			devAssignment := payload.DeveloperTaskAssignment{
				Developer: dev,
				Tasks:     []payload.Task{},
			}

			var updatedTasks []payload.Task
			for _, task := range tasks {
				if task.Duration <= remainingHours {
					devAssignment.Tasks = append(devAssignment.Tasks, task)
					remainingHours -= task.Duration
					developerWorkloads[dev.ID] += task.Duration
				} else {
					updatedTasks = append(updatedTasks, task)
				}
			}
			tasks = updatedTasks

			if len(devAssignment.Tasks) > 0 {
				assignment.DeveloperTasks = append(assignment.DeveloperTasks, devAssignment)
			}
		}

		if len(assignment.DeveloperTasks) > 0 {
			assignments = append(assignments, assignment)
			totalWeeks++
		}
	}

	resp := payload.ScheduleAssignmentResponse{
		Assignments: assignments,
		MinDuration: totalWeeks,
	}

	s.logger.Trace("Assignments scheduled successfully with minDuration=%v weeks", totalWeeks)
	return resp, nil
}
