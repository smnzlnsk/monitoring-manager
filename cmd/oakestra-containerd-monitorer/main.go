package main

import (
	"context"
	"fmt"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/errdefs"
	"google.golang.org/grpc"
)

func main() {
	// Connect to containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")

	// List all running containers
	containers, err := client.Containers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		// task, err := container.Task(ctx, nil)
		if err != nil {
			if errdefs.IsNotFound(err) {
				continue
			}
			log.Fatal(err)
		}

		// Get metrics for the task
		conn, err := grpc.NewClient(client.Conn().Target(), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		taskClient := tasks.NewTasksClient(conn)
		metricsResponse, err := taskClient.Metrics(ctx, &tasks.MetricsRequest{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Metrics for container %s:\n", container.ID())
		for _, metric := range metricsResponse.Metrics {
			fmt.Printf("  %s: %v\n", metric.ID, metric)
		}
	}

	fmt.Println("Metrics extracted successfully")
}
