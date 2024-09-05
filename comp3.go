package main

import (
	"context"
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp3 struct {
	bg.BasicComponent
}

func (c *Comp3) ProcessReq(ctx context.Context, request bg.CompRequest[interface{}]) {
	//fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Println("Component ", c.CompId, " processing on ", c.GetStagingVersion())
}

func (c *Comp3) Sync(ctx context.Context) {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp3) CancelReq(ctx context.Context) {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
