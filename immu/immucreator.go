package immu

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/context"
	"log"
	"time"
)

const (
	ContainerUser     = "immudb"
	ContainerPassword = "immudb"
	ContainerDb       = "defaultdb"
)

type TestImmu struct {
	container testcontainers.Container
	Port      nat.Port
}

func SetupTestImmu() *TestImmu {

	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, port, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	cancel()

	return &TestImmu{
		container: container,
		Port:      port,
	}
}

func createContainer(ctx context.Context) (testcontainers.Container, nat.Port, error) {
	var port = "3322/tcp"
	var port2 = "9497/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "codenotary/immudb:latest",
			ExposedPorts: []string{port, port2},
			WaitingFor: wait.ForAll(
				wait.ForLog("sessions guard started"),
				wait.ForListeningPort("3322/tcp"),
				wait.ForListeningPort("9497/tcp"),
			),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)

	if err != nil {
		return container, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "3322")
	if err != nil {
		return container, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("immudb container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	return container, p, nil
}
