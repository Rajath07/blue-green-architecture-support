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
	components := bg.InitializeComponents(ctx, "dependency.yaml", compCollec, 2)
	components.SendReq("Comp1", bg.Update, 10, 0)
	components.SendReq("Comp1", bg.Create, 33, 2)
	components.SendReq("Comp1", bg.Create, 60, 1)
	//components.SendReq("Comp4", bg.Create, 34, 0)
	time.Sleep(20 * time.Second)

	//supervisor.GetChannel(1) <- "Start Processing"
	//cancel()
	//components[1].ProcessReq(ctx)
}
