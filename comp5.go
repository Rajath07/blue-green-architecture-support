package main

import (
	"fmt"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp5 struct {
	bg.BasicComponent
}

func (c *Comp5) ProcessReq(request bg.CompRequest[interface{}]) {
	//fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Println("Component ", c.CompId, " processing on ", c.GetStagingVersion())
	// fmt.Println("Sleeping for 5 seconds")
	// time.Sleep(5 * time.Second)
}

func (c *Comp5) Sync() {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp5) CancelReq() {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
