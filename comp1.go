package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp1 struct {
	bg.BasicComponent
}

func (c *Comp1) ProcessReq(ctx context.Context) {
	//c.OutChannel[1] <- "Start Processing"
	//fmt.Println("Outchannels: ", c.OutChannel)
	fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Printf("\nComp%d processing from userdefined\n", c.CompId)
	fmt.Println("Sleeping for 2 seconds")
	time.Sleep(2 * time.Second)
}

func (c *Comp1) Switch(ctx context.Context) {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp1) CancelReq(ctx context.Context) {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
