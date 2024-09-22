package main

import (
	"fmt"
	"time"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp2 struct {
	bg.BasicComponent
}

func (c *Comp2) ProcessReq(request bg.CompRequest[interface{}]) {
	// c.OutChannel[0] <- "Start Processing"
	// c.OutChannel[1] <- "Start Processing"
	//fmt.Println("Printing my customComp ID ", c.CompId)
	fmt.Println("Component ", c.CompId, " processing on ", c.GetStagingVersion())
	time.Sleep(1 * time.Second)
}

func (c *Comp2) Sync() {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp2) Cancel() {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
