package main

import (
	"context"
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp2 struct {
	bg.BasicComponent
}

func (c *Comp2) ProcessReq(ctx context.Context) {
	// c.OutChannel[0] <- "Start Processing"
	// c.OutChannel[1] <- "Start Processing"
	fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Printf("\nComp%d processing from userdefined\n", c.CompId)
}

func (c *Comp2) Switch(ctx context.Context) {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp2) CancelReq(ctx context.Context) {
	fmt.Printf("Comp1 cancelling from userdefined")
}
