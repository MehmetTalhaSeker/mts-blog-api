package database_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/database"
)

const (
	user     = "test-user"
	name     = "test-name"
	password = "test-password"
)

// TestNewPostgresStore tests the NewPostgresStore function.
func TestNewPostgresStore(t *testing.T) {
	ctx := context.Background()

	port := "5432/tcp"
	env := map[string]string{
		"POSTGRES_PASSWORD": password,
		"POSTGRES_USER":     user,
		"POSTGRES_DB":       name,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Clean up the container.
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("container termination failed: %v", err)
		}
	}()

	// Check store is nil or not.
	store := database.NewPostgresStore(database.WithUser(user), database.WithName(name), database.WithPassword(password), database.WithPort(p.Port()))

	if store == nil {
		t.Fatalf("NewPostgresStore returned a nil store")
	}

	defer store.GetInstance().Close()
}
