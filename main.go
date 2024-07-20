package main

import (
	"context"
	"fmt"

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
	fmt.Println("Components initialized")

	components["ComponentA"].ProcessReq(ctx)

	supervisor.GetChannel("ComponentA") <- "Start Processing"
	//cancel()
	components["ComponentA"].ProcessReq(ctx)

	components["ComponentA"].ProcessReq(ctx)
}
