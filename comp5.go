package main

import (
	"context"
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp5 struct {
	bg.BasicComponent
}

func (c *Comp5) ProcessReq(ctx context.Context) {
	fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Printf("\nComp%d processing from userdefined\n", c.CompId)
	// fmt.Println("Sleeping for 5 seconds")
	// time.Sleep(5 * time.Second)
}

func (c *Comp5) Switch(ctx context.Context) {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp5) CancelReq(ctx context.Context) {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
