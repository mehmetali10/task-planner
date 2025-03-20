package payload

import "time"

type (
	Task struct {
		ID         uint       `json:"id"`
		ExternalID uint       `json:"externalId"`
		Name       string     `json:"name"`
		Duration   int        `json:"duration"`
		Difficulty int        `json:"difficulty"`
		Provider   string     `json:"provider"`
		CreatedAt  *time.Time `json:"createdAt"`
		UpdatedAt  *time.Time `json:"updatedAt"`
	}

	CreateTaskRequest struct {
		ExternalID uint   `json:"externalId" validate:"required"`
		Name       string `json:"name" validate:"required,min=3,max=100"`
		Duration   int    `json:"duration" validate:"required,min=1,max=1000"`
		Difficulty int    `json:"difficulty" validate:"required,min=1,max=10"`
		Provider   string `json:"provider" validate:"required,min=3,max=150"`
	}
	CreateTaskResponse struct {
		ID        uint       `json:"id"`
		CreatedAt *time.Time `json:"createdAt"`
	}

	ListTasksRequest struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
	ListTasksResponse struct {
		Tasks []Task `json:"tasks"`
	}
)

type (
	Assignment struct {
		DeveloperTasks []DeveloperTaskAssignment `json:"developerTasks"`
	}

	DeveloperTaskAssignment struct {
		Developer Developer `json:"developer"`
		Tasks     []Task    `json:"tasks"`
	}

	ScheduleAssignmentRequest struct {
	}

	ScheduleAssignmentResponse struct {
		Assignments []Assignment `json:"assignments"`
		MinDuration int          `json:"minDuration"`
	}
)

type (
	Developer struct {
		ID        uint       `json:"id"`
		FirstName string     `json:"firstName"`
		LastName  string     `json:"lastName"`
		Email     string     `json:"email"`
		CreatedAt *time.Time `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	ListDevelopersRequest  struct{}
	ListDevelopersResponse struct {
		Developers []Developer `json:"developers"`
	}
)
