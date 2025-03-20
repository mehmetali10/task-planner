package service

import (
	"context"
	"sort"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
)

func (s *service) ScheduleAssignments(ctx context.Context, req payload.ScheduleAssignmentRequest) (payload.ScheduleAssignmentResponse, error) {
	s.logger.Trace("Scheduling assignments")

	// Fetch the list of tasks and developers
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

	if len(tasks) == 0 || len(developers) == 0 {
		s.logger.Warn("No tasks or developers available")
		return payload.ScheduleAssignmentResponse{}, nil
	}

	const weeklyWorkHours = 45
	const workDaysInWeek = 5 // 5 workdays in a week
	assignments := []payload.Assignment{}
	totalWeeks := 0
	totalElapsedWorkHour := 0.0

	// Sort tasks in descending order of duration
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Duration > tasks[j].Duration
	})

	// Assignment process
	remainingTasks := tasks

	for len(remainingTasks) > 0 {
		assignment := payload.Assignment{
			DeveloperTasks: []payload.DeveloperTaskAssignment{},
		}

		// Reset weekly workloads for developers
		developerWorkloads := make(map[uint]float64)
		for _, dev := range developers {
			developerWorkloads[dev.ID] = 0
		}

		// Determine tasks to be assigned for the week
		newRemainingTasks := []payload.Task{}

		for _, task := range remainingTasks {
			// Find the developer with the least workload
			sort.Slice(developers, func(i, j int) bool {
				return developerWorkloads[developers[i].ID] < developerWorkloads[developers[j].ID]
			})

			assigned := false
			for i := range developers {
				dev := &developers[i]

				// The amount of difficulty a developer can handle in 1 hour
				if dev.Capacity <= 0 {
					continue
				}

				effectiveTaskDuration := float64(task.Difficulty) / float64(dev.Capacity)

				if developerWorkloads[dev.ID]+effectiveTaskDuration <= float64(weeklyWorkHours) {
					// Assign task to developer
					found := false
					for j := range assignment.DeveloperTasks {
						if assignment.DeveloperTasks[j].Developer.ID == dev.ID {
							assignment.DeveloperTasks[j].Tasks = append(assignment.DeveloperTasks[j].Tasks, task)
							found = true
							break
						}
					}
					if !found {
						assignment.DeveloperTasks = append(assignment.DeveloperTasks, payload.DeveloperTaskAssignment{
							Developer: *dev,
							Tasks:     []payload.Task{task},
						})
					}
					developerWorkloads[dev.ID] += effectiveTaskDuration
					totalElapsedWorkHour += effectiveTaskDuration
					assigned = true
					break
				}
			}

			// If the task cannot be assigned, carry it over to the next week
			if !assigned {
				newRemainingTasks = append(newRemainingTasks, task)
			}
		}

		// Check and adjust workload balance
		sort.Slice(assignment.DeveloperTasks, func(i, j int) bool {
			return len(assignment.DeveloperTasks[i].Tasks) > len(assignment.DeveloperTasks[j].Tasks)
		})

		for i := 0; i < len(assignment.DeveloperTasks); i++ {
			for j := 0; j < len(assignment.DeveloperTasks); j++ {
				if len(assignment.DeveloperTasks[i].Tasks) > 1 && len(assignment.DeveloperTasks[j].Tasks) == 0 {
					// Move a task from an overloaded developer to an idle developer
					taskToMove := assignment.DeveloperTasks[i].Tasks[len(assignment.DeveloperTasks[i].Tasks)-1]
					assignment.DeveloperTasks[j].Tasks = append(assignment.DeveloperTasks[j].Tasks, taskToMove)
					assignment.DeveloperTasks[i].Tasks = assignment.DeveloperTasks[i].Tasks[:len(assignment.DeveloperTasks[i].Tasks)-1]

					// Update workloads
					developerWorkloads[assignment.DeveloperTasks[j].Developer.ID] += float64(taskToMove.Difficulty) / float64(assignment.DeveloperTasks[j].Developer.Capacity)
					developerWorkloads[assignment.DeveloperTasks[i].Developer.ID] -= float64(taskToMove.Difficulty) / float64(assignment.DeveloperTasks[i].Developer.Capacity)
					break
				}
			}
		}

		// If tasks are assigned, move to the next week
		if len(assignment.DeveloperTasks) > 0 {
			assignments = append(assignments, assignment)
			totalWeeks++

			// Calculate and log the minimum duration in days
			totalDays := totalWeeks * workDaysInWeek
			s.logger.Trace("Week %d completed, total duration so far: %d days", totalWeeks, totalDays)
		}

		// Update remaining tasks
		remainingTasks = newRemainingTasks
	}

	// Calculate the minimum duration in both weeks and days
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
