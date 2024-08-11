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
	fmt.Println("Comp2 processing from userdefined")
}

func (c *Comp2) SyncReq(ctx context.Context) {
	fmt.Printf("Comp1 syncing from userdefined")
}

func (c *Comp2) CancelReq(ctx context.Context) {
	fmt.Printf("Comp1 cancelling from userdefined")
}
