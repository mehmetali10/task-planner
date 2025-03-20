package service

import (
	"context"
	"sort"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
)

func (s *service) ScheduleAssignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	s.logger.Trace("Scheduling assignments")

	// Fetch the list of tasks and developers
	tasks, err := s.fetchTasks(ctx)
	if err != nil {
		return payload.ScheduleAssignmentResponse{}, err
	}

	developers, err := s.fetchDevelopers(ctx)
	if err != nil {
		return payload.ScheduleAssignmentResponse{}, err
	}

	// Return empty response if no tasks or developers
	if len(tasks) == 0 || len(developers) == 0 {
		s.logger.Warn("No tasks or developers available")
		return payload.ScheduleAssignmentResponse{}, nil
	}

	// Constants for weekly work hours and days
	const weeklyWorkHours = 4
	const workDaysInWeek = 5 // 5 workdays per week
	assignments := []payload.Assignment{}
	totalWeeks := 0
	totalElapsedWorkHour := 0.0

	// Sort tasks by duration in descending order
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Duration > tasks[j].Duration
	})

	remainingTasks := tasks
	for len(remainingTasks) > 0 {
		// Create an assignment for the week
		var assignment payload.Assignment
		var newRemainingTasks []payload.Task

		assignment, newRemainingTasks, totalElapsedWorkHour = s.assignTasksToDevelopers(remainingTasks, developers, weeklyWorkHours, totalElapsedWorkHour)

		// If tasks are assigned, process the result
		if len(assignment.DeveloperTasks) > 0 {
			assignments = append(assignments, assignment)
			totalWeeks++

			// Log weekly progress
			totalDays := totalWeeks * workDaysInWeek
			s.logger.Trace("Week %d completed, total duration so far: %d days", totalWeeks, totalDays)
		}

		// Update remaining tasks for the next iteration
		remainingTasks = newRemainingTasks
	}

	// Calculate the minimum total days and hours
	minDays := totalWeeks * workDaysInWeek
	resp := payload.ScheduleAssignmentResponse{
		Assignments:          assignments,
		MinWeek:              uint(totalWeeks),
		TotalWorkDay:         uint(minDays),
		TotalElapsedWorkHour: uint(totalElapsedWorkHour),
	}

	s.logger.Trace("Assignments scheduled successfully with minWeek=%v weeks (%v days), totalElapsedWorkHour=%v hours", totalWeeks, minDays, totalElapsedWorkHour)
	return resp, nil
}

// fetchTasks retrieves the list of tasks from the repository.
func (s *service) fetchTasks(ctx context.Context) ([]payload.Task, error) {
	tasksResp, err := s.repository.ListTasks(ctx, payload.ListTasksRequest{})
	if err != nil {
		s.logger.Error("Failed to list tasks: error=%v", err)
		return nil, err
	}
	return tasksResp.Tasks, nil
}

// fetchDevelopers retrieves the list of developers from the repository.
func (s *service) fetchDevelopers(ctx context.Context) ([]payload.Developer, error) {
	developersResp, err := s.repository.ListDevelopers(ctx, payload.ListDevelopersRequest{})
	if err != nil {
		s.logger.Error("Failed to list developers: error=%v", err)
		return nil, err
	}
	return developersResp.Developers, nil
}

// assignTasksToDevelopers assigns tasks to developers based on their capacity and workload.
func (s *service) assignTasksToDevelopers(remainingTasks []payload.Task, developers []payload.Developer, weeklyWorkHours int, totalElapsedWorkHour float64) (payload.Assignment, []payload.Task, float64) {
	assignment := payload.Assignment{
		DeveloperTasks: []payload.DeveloperTaskAssignment{},
	}
	newRemainingTasks := []payload.Task{}

	// Initialize developer workloads
	developerWorkloads := make(map[uint]float64)
	for _, dev := range developers {
		developerWorkloads[dev.ID] = 0
	}

	// Assign tasks to developers
	for _, task := range remainingTasks {
		assigned := false
		for _, dev := range developers {
			// Calculate the effective task duration based on developer's capacity
			effectiveTaskDuration := float64(task.Difficulty) / float64(dev.Capacity)

			// Check if the developer can handle the task within weekly limits
			if developerWorkloads[dev.ID]+effectiveTaskDuration <= float64(weeklyWorkHours) {
				// Assign task to developer
				assignment = s.assignTaskToDeveloper(dev, task, assignment)
				developerWorkloads[dev.ID] += effectiveTaskDuration
				totalElapsedWorkHour += effectiveTaskDuration
				assigned = true
				break
			}
		}

		// If the task could not be assigned, add it to the remaining tasks for the next week
		if !assigned {
			newRemainingTasks = append(newRemainingTasks, task)
		}
	}

	// Rebalance tasks between developers to ensure fair workload distribution
	s.rebalanceWorkload(assignment, developerWorkloads)

	return assignment, newRemainingTasks, totalElapsedWorkHour
}

// assignTaskToDeveloper assigns a task to a developer and returns the updated assignment.
func (s *service) assignTaskToDeveloper(dev payload.Developer, task payload.Task, assignment payload.Assignment) payload.Assignment {
	// Check if the task is already assigned to the developer
	for i := range assignment.DeveloperTasks {
		if assignment.DeveloperTasks[i].Developer.ID == dev.ID {
			assignment.DeveloperTasks[i].Tasks = append(assignment.DeveloperTasks[i].Tasks, task)
			return assignment
		}
	}

	// If not, create a new assignment for the developer
	assignment.DeveloperTasks = append(assignment.DeveloperTasks, payload.DeveloperTaskAssignment{
		Developer: dev,
		Tasks:     []payload.Task{task},
	})
	return assignment
}

// rebalanceWorkload ensures the tasks are fairly distributed among developers.
func (s *service) rebalanceWorkload(assignment payload.Assignment, developerWorkloads map[uint]float64) {
	// Sort developer tasks by the number of tasks assigned
	sort.Slice(assignment.DeveloperTasks, func(i, j int) bool {
		return len(assignment.DeveloperTasks[i].Tasks) > len(assignment.DeveloperTasks[j].Tasks)
	})

	// Rebalance tasks to developers with less workload
	for i := 0; i < len(assignment.DeveloperTasks); i++ {
		for j := 0; j < len(assignment.DeveloperTasks); j++ {
			if len(assignment.DeveloperTasks[i].Tasks) > 1 && len(assignment.DeveloperTasks[j].Tasks) == 0 {
				// Move a task from an overloaded developer to an underloaded one
				taskToMove := assignment.DeveloperTasks[i].Tasks[len(assignment.DeveloperTasks[i].Tasks)-1]
				assignment.DeveloperTasks[j].Tasks = append(assignment.DeveloperTasks[j].Tasks, taskToMove)
				assignment.DeveloperTasks[i].Tasks = assignment.DeveloperTasks[i].Tasks[:len(assignment.DeveloperTasks[i].Tasks)-1]

				// Update workload of both developers
				developerWorkloads[assignment.DeveloperTasks[j].Developer.ID] += float64(taskToMove.Difficulty) / float64(assignment.DeveloperTasks[j].Developer.Capacity)
				developerWorkloads[assignment.DeveloperTasks[i].Developer.ID] -= float64(taskToMove.Difficulty) / float64(assignment.DeveloperTasks[i].Developer.Capacity)
				break
			}
		}
	}
}
