package main

import (
	"context"
	"time"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

func main() {
	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define component names and dependencies
	names := []string{"ComponentA", "ComponentB"}
	dependencies := []bg.Dependency{
		{Child: "ComponentB", Parent: "ComponentA"},
	}

	// Create the supervisor
	supervisor := bg.NewSupervisor(names)

	// Initialize and run components
	components := bg.InitializeComponents(ctx, supervisor, names, dependencies)

	// Simulate sending a signal to each component's supervisor channel
	time.Sleep(1 * time.Second)
	for _, name := range names {
		supervisor.GetChannel(name) <- "Start Processing"
	}

	// Simulate a cancel signal after some time
	time.Sleep(2 * time.Second)
	cancel()

	// Wait for all components to finish
	// (This is handled within the InitializeComponents function)
}
