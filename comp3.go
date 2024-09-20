package main

import (
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp3 struct {
	bg.BasicComponent
}

func (c *Comp3) ProcessReq(request bg.CompRequest[interface{}]) {
	//fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Println("Component ", c.CompId, " processing on ", c.GetStagingVersion())
}

func (c *Comp3) Sync() {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp3) CancelReq() {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
