package main

import (
	"context"
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp1 struct {
	bg.BasicComponent
}

func (c *Comp1) ProcessReq(ctx context.Context) {
	//c.OutChannel[0] <- "Start Processing"
	fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Printf("\nComp%d processing from userdefined\n", c.CompId)
}

func (c *Comp1) SyncReq(ctx context.Context) {
	fmt.Printf("Comp1 syncing from userdefined")
}

func (c *Comp1) CancelReq(ctx context.Context) {
	fmt.Printf("Comp1 cancelling from userdefined")
}
