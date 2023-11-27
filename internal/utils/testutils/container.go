package testutils

import (
	"context"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresTestContainer(ctx context.Context) (testcontainers.Container, nat.Port) {
	port := "5432/tcp"
	env := map[string]string{
		"POSTGRES_PASSWORD": "test-password",
		"POSTGRES_USER":     "test-user",
		"POSTGRES_DB":       "test-name",
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
		log.Fatalf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return container, p
}

func TerminateContainer(ctx context.Context, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("container.Terminate(%v): %v", ctx, err)
	}

	log.Println("container.Terminate() success")
}
