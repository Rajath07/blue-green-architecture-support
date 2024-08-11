package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

func main() {
	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the supervisor
	// supervisor := bg.NewSupervisor(compIds)
	// supervisor.
	compCollec := []bg.Component{}
	customComp1 := &Comp1{}
	customComp2 := &Comp2{}
	customComp3 := &Comp3{}
	customComp4 := &Comp4{}
	customComp5 := &Comp5{}
	compCollec = append(compCollec, customComp1)
	compCollec = append(compCollec, customComp2)
	compCollec = append(compCollec, customComp3)
	compCollec = append(compCollec, customComp4)
	compCollec = append(compCollec, customComp5)

	// Initialize and run components
	components := bg.InitializeComponents(ctx, "dependency.yaml", compCollec)
	fmt.Println("Components initialized ")
	components.SendReq("Comp1", bg.Create, 10)
	//components["Comp2"].ProcessReq(ctx)

	//components[1].ProcessReq(ctx)
	// components[2].ProcessReq(ctx)
	//components[10].ProcessReq(ctx)
	time.Sleep(5 * time.Second)

	//supervisor.GetChannel(1) <- "Start Processing"
	//cancel()
	//components[1].ProcessReq(ctx)
}
