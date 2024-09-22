package main

import (
	"fmt"
	"time"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type Comp1 struct {
	bg.BasicComponent
	blue  []int
	green []int
}

func (c *Comp1) ProcessReq(request bg.Request[interface{}]) {
	//c.OutChannel[1] <- "Start Processing"
	//fmt.Println("Outchannels: ", c.OutChannel)
	//fmt.Println("Printing my customComp ID ", c.CompId)
	//fmt.Printf("\nComp%d processing from userdefined\n", c.CompId)
	fmt.Println("Comp ", c.CompId, " received request in user defined", request)
	fmt.Println("Component ", c.CompId, " processing on ", c.GetStagingVersion())
	time.Sleep(2 * time.Second)
}

func (c *Comp1) Sync() {
	fmt.Println("Switching blue to green in ", c.CompId)
}

func (c *Comp1) Cancel() {
	fmt.Printf("Component %d cancelling from userdefined", c.CompId)
}
