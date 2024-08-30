package main

import (
	"context"
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp3 struct {
	bg.BasicComponent
}

func (c *Comp3) ProcessReq(ctx context.Context) {
	fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Printf("\nComp%d processing from userdefined\n", c.CompId)
}

func (c *Comp3) Switch(ctx context.Context) {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp3) CancelReq(ctx context.Context) {
	fmt.Printf("Comp3 cancelling from userdefined")
}
