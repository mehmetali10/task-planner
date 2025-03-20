package postgres_repository

import (
	"context"
	"testing"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/pkg/testcontainer"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	// Start the PostgreSQL container and perform migrations once before the tests
	cleanup := testcontainer.StartPostgresContainer(t)
	defer cleanup()

	// Initialize the repository
	repo := NewPostgresRepo()

	t.Run("CreateTask", func(t *testing.T) {
		tests := []struct {
			name          string
			input         payload.CreateTaskRequest
			expectedError bool
		}{
			{
				name: "CreateTask_Success",
				input: payload.CreateTaskRequest{
					ExternalID: 123,
					Name:       "New Task",
					Duration:   5,
					Difficulty: 3,
					Provider:   "Test Provider",
				},
				expectedError: false,
			},
			{
				name: "CreateTask_TaskAlreadyExists",
				input: payload.CreateTaskRequest{
					ExternalID: 123,
					Name:       "New Task",
					Duration:   5,
					Difficulty: 3,
					Provider:   "Test Provider",
				},
				expectedError: true,
			},
		}

		// Run each test case
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := repo.CreateTask(context.Background(), tt.input)

				// Handle errors and validate response
				if tt.expectedError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Validate that the response ID is greater than 0
					require.Greater(t, resp.ID, uint(0))
				}
			})
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		tests := []struct {
			name          string
			input         payload.ListTasksRequest
			expectedError bool
		}{
			{
				name: "ListTasks_Success",
				input: payload.ListTasksRequest{
					Limit:  10,
					Offset: 0,
				},
				expectedError: false,
			},
			{
				name: "ListTasks_NoTasksFound",
				input: payload.ListTasksRequest{
					Limit:  0,
					Offset: 0,
				},
				expectedError: false,
			},
		}

		// Run each test case
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := repo.ListTasks(context.Background(), tt.input)

				// Handle errors and validate response
				if tt.expectedError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Validate that we got some tasks back, if expected
					require.NotNil(t, resp.Tasks)
					require.Greater(t, len(resp.Tasks), 0) // Expecting at least 1 task
				}
			})
		}
	})

	t.Run("ListDevelopers", func(t *testing.T) {
		tests := []struct {
			name          string
			input         payload.ListDevelopersRequest
			expectedError bool
		}{
			{
				name:          "ListDevelopers_Success",
				input:         payload.ListDevelopersRequest{},
				expectedError: false,
			},
		}

		// Run each test case
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := repo.ListDevelopers(context.Background(), tt.input)

				// Handle errors and validate response
				if tt.expectedError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Validate that we got some developers back, if expected
					require.NotNil(t, resp.Developers)
					require.Greater(t, len(resp.Developers), 0) // Expecting at least 1 developer
				}
			})
		}
	})
}
