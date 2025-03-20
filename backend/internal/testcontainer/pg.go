package testcontainer

import (
	"context"
	"testing"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"github.com/mehmetali10/task-planner/internal/task/migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// StartPostgresContainer starts a PostgreSQL container for testing purposes
func StartPostgresContainer(t *testing.T) func() {
	ctx := context.Background()
	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("task"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("pass"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	require.NoError(t, err)

	host, err := ctr.Host(ctx)
	require.NoError(t, err)

	port, err := ctr.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Configure the app's database settings
	config.LoadConfig()
	config.GetApp().DBHost = host
	config.GetApp().DBPort = port.Int()
	config.GetApp().DBName = "task"
	config.GetApp().DBUser = "postgres"
	config.GetApp().DBPassword = "pass"

	// Run migrations
	migrate.MigrateAndSeed()
	// Return the host, port, and cleanup function
	return func() {
		// Cleanup container
		testcontainers.CleanupContainer(t, ctr)
	}
}
