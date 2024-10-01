package main

import (
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp4 struct {
	bg.BasicComponent
}

func (c *Comp4) ProcessReq(request bg.Request[interface{}]) {
	//fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Println("Component ", c.CompId, " processing on ", c.GetStagingVersion())
}

func (c *Comp4) Sync() {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp4) Cancel() {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
