package main

import (
	"context"
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp4 struct {
	bg.BasicComponent
}

func (c *Comp4) ProcessReq(ctx context.Context) {
	fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Println("Comp3 processing from userdefined")
}

func (c *Comp4) SyncReq(ctx context.Context) {
	fmt.Printf("Comp3 syncing from userdefined")
}

func (c *Comp4) CancelReq(ctx context.Context) {
	fmt.Printf("Comp3 cancelling from userdefined")
}
