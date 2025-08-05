package postgres

import (
	"context"
	"strings"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap/zaptest"
)

const (
	dbName     = "users"
	dbUser     = "user"
	dbPassword = "password"
)

func terminateContainer(
	t *testing.T, container testcontainers.Container, pg *repository.Repository,
) {
	t.Helper()

	assert.NoError(t, pg.DB.Close(), "failed close db")
	require.NoError(
		t,
		testcontainers.TerminateContainer(container),
		"failed to terminate container",
	)
}

func startPostgresContainer(t *testing.T) *postgres.PostgresContainer {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		require.NoError(t, err, "failed to start container")

		return nil
	}

	return postgresContainer
}

func getPostgres(t *testing.T) (*repository.Repository, *postgres.PostgresContainer) {
	t.Helper()

	container := startPostgresContainer(t)
	require.NotNil(t, container)

	ctx := context.Background()

	endpoint, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	endpointSplit := strings.Split(endpoint, ":")
	cfg := config.Postgres{
		User:     dbUser,
		Password: dbPassword,
		Host:     endpointSplit[0],
		Port:     endpointSplit[1],
		DB:       dbName,
		MaxConns: "20",
	}

	log := zaptest.NewLogger(t)
	pg, err := New(ctx, &cfg, log)
	require.NoError(t, err)

	return pg, container
}
