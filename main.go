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
	// data, err := bg.ParseYAML("dependency.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(data)

	// Create the supervisor
	//supervisor := bg.NewSupervisor(compIds)
	compCollec := []bg.Component{}
	customComp1 := &Comp1{}
	customComp2 := &Comp2{}
	compCollec = append(compCollec, customComp1)
	compCollec = append(compCollec, customComp2)

	// Initialize and run components
	components := bg.InitializeComponents(ctx, "dependency.yaml", compCollec)
	fmt.Println("Components initialized ", components)
	components["Comp1"].ProcessReq(ctx)

	//components[1].ProcessReq(ctx)
	// components[2].ProcessReq(ctx)
	//components[10].ProcessReq(ctx)
	time.Sleep(5 * time.Second)

	//supervisor.GetChannel(1) <- "Start Processing"
	//cancel()
	//components[1].ProcessReq(ctx)
}
