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
	compIds := []int{1, 2}
	dependencies := []bg.Dependency{
		{Child: 2, Parent: 1},
	}

	// Create the supervisor
	supervisor := bg.NewSupervisor(compIds)

	// Initialize and run components
	components := bg.InitializeComponents(ctx, supervisor, compIds, dependencies)
	fmt.Println("Components initialized")

	components[1].ProcessReq(ctx)

	supervisor.GetChannel(1) <- "Start Processing"
	//cancel()
	components[1].ProcessReq(ctx)

	components[1].ProcessReq(ctx)
}
