package storage_test

import (
	"context"
	"database/sql"
	"github.com/serjnester/weather-subscription-service/domain/models"
	"github.com/serjnester/weather-subscription-service/internal/storage"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "weather_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(10 * time.Second),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := "postgres://postgres:secret@" + host + ":" + port.Port() + "/weather_test?sslmode=disable"

	var db *sql.DB
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil && db.Ping() == nil {
			break
		}
		time.Sleep(time.Second)
	}
	require.NoError(t, err)

	// Apply schema
	schema, err := os.ReadFile("testdata/schema.sql")
	require.NoError(t, err)

	_, err = db.Exec(string(schema))
	require.NoError(t, err)

	cleanup := func() {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestSubscriptionStorage(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	subStorage := storage.NewSubscriptionStorage(storage.New(db))

	ctx := context.Background()
	sub := models.Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Token:     "abc123",
		Frequency: "daily",
	}

	// Insert
	err := subStorage.Create(ctx, sub)
	require.NoError(t, err)

	// GetByToken
	got, err := subStorage.GetByToken(ctx, sub.Token)
	require.NoError(t, err)
	require.Equal(t, sub.Email, got.Email)

	// Confirm
	err = subStorage.Confirm(ctx, sub.Token)
	require.NoError(t, err)

	confirmed, err := subStorage.GetByToken(ctx, sub.Token)
	require.NoError(t, err)
	require.True(t, confirmed.Confirmed)

	// Unsubscribe
	err = subStorage.Unsubscribe(ctx, sub.Token)
	require.NoError(t, err)

	_, err = subStorage.GetByToken(ctx, sub.Token)
	require.Error(t, err)
}
