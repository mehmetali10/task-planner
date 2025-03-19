package payload

type (
	Task struct {
		ID         uint64 `json:"id"`
		ExternalID uint64 `json:"externalId"`
		Name       string `json:"name"`
		Duration   int    `json:"duration"`
		Difficulty int    `json:"difficulty"`
		Provider   string `json:"provider"`
		CreatedAt  string `json:"createdAt"`
		UpdatedAt  string `json:"updatedAt"`
	}

	CreateTaskRequest struct {
		ExternalID uint   `json:"externalId" validate:"required"`
		Name       string `json:"name" validate:"required,min=3,max=100"`
		Duration   int    `json:"duration" validate:"required,min=1,max=1000"`
		Difficulty int    `json:"difficulty" validate:"required,min=1,max=10"`
		Provider   string `json:"provider" validate:"required,min=3,max=150"`
	}
	CreateTaskResponse struct {
		ID        uint   `json:"id"`
		CreatedAt string `json:"createdAt"`
	}

	ListTasksRequest struct {
		Offset int `json:"offset" validate:"min=0"`
		Limit  int `json:"limit" validate:"min=1,max=100"`
	}
	ListTasksResponse struct {
		Tasks []Task `json:"tasks"`
	}
)

type (
	Assignment struct {
		ID          uint      `json:"id"`
		TaskID      uint      `json:"taskId"`
		DeveloperID uint      `json:"developerId"`
		Task        Task      `json:"task"`
		Developer   Developer `json:"developer"`
		CreatedAt   string    `json:"createdAt"`
		UpdatedAt   string    `json:"updatedAt"`
	}

	ScheduleAssignmentRequest struct {
	}
	ScheduleAssignmentResponse struct {
		Assignments []Assignment `json:"assignments"`
	}

	ListAssignmentsRequest  struct{}
	ListAssignmentsResponse struct{}
)

type (
	Developer struct {
		ID        uint   `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	ListDevelopersRequest  struct{}
	ListDevelopersResponse struct {
		Developers []Developer `json:"developers"`
	}
)
