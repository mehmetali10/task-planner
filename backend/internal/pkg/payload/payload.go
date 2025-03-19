package payload

type (
	Task struct {
		ID         uint64 `json:"id"`
		ExternalID uint64 `json:"external_id"`
		Name       string `json:"name"`
		Duration   int    `json:"duration"`
		Difficulty int    `json:"difficulty"`
		Provider   string `json:"provider"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	CreateTaskRequest struct {
		ExternalID uint64 `json:"external_id"`
		Name       string `json:"name"`
		Duration   int    `json:"duration"`
		Difficulty int    `json:"difficulty"`
		Provider   string `json:"provider"`
	}
	CreateTaskResponse struct {
		ID        uint64 `json:"id"`
		CreatedAt string `json:"created_at"`
	}

	GetTaskRequest struct {
		ID uint64 `json:"id"`
	}
	GetTaskResponse struct {
		Task Task `json:"task"`
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
		ID          uint64 `json:"id"`
		TaskID      uint64 `json:"task_id"`
		DeveloperID uint64 `json:"developer_id"`
		Week        int    `json:"week"`
		Task        Task   `json:"task"`
		Developer   Developer
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	ScheduleAssignmentRequest struct {
		Week int `json:"week"`
	}
	ScheduleAssignmentResponse struct {
		Assignments []Assignment `json:"assignments"`
	}

	ListAssignmentsRequest  struct{}
	ListAssignmentsResponse struct{}
)

type (
	Developer struct {
		ID        uint64 `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	ListDevelopersRequest  struct{}
	ListDevelopersResponse struct {
		Developers []Developer `json:"developers"`
	}
)
