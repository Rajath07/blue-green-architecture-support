package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type CompA struct {
	bg.BasicComponent
	arrBlue   []int
	arrGreen  []int
	permBlue  []int
	permGreen []int
}

// Function to calculate the permutation of an array
func (c *CompA) calculatePerm(arr []int) []int {
	perm := make([]int, len(arr))
	for i := range perm {
		perm[i] = i
	}

	// Sort the permutation array based on the corresponding values in arr
	sort.Slice(perm, func(i, j int) bool {
		return arr[perm[i]] < arr[perm[j]]
	})
	return perm
}

func (c *CompA) ProcessReq(request bg.Request[interface{}]) {
	stagingVersion := c.GetStagingVersion()

	if request.ComponentName == reflect.TypeOf(c).Elem().Name() {
		if request.Operation == bg.Update {
			switch stagingVersion {
			case 0:
				c.arrBlue[request.Index] = request.Data.(int)
				//fmt.Println("Appended to blue data")
			case 1:
				c.arrGreen[request.Index] = request.Data.(int)
				//fmt.Println("Appended to green data")
			default:
				fmt.Println("Unknown staging version:", stagingVersion)

			}

		}
	}
	switch stagingVersion {
	case 0:
		c.permBlue = c.calculatePerm(c.arrBlue)
	case 1:
		c.permGreen = c.calculatePerm(c.arrGreen)
	default:
		fmt.Println("Unknown staging version:", stagingVersion)

	}
	fmt.Println("Process completed in CompA")
}

func (c *CompA) Sync() {
	stagingVersion := c.GetStagingVersion()

	switch stagingVersion {
	case 0:
		c.arrBlue = c.arrGreen
		c.permBlue = c.permGreen
	case 1:
		c.arrGreen = c.arrBlue
		c.permGreen = c.permBlue
	default:
		fmt.Println("Unknown staging version:", stagingVersion)
	}
	fmt.Println("Sync completed in CompA")

}

func (c *CompA) Cancel() {
	fmt.Printf("Component A cancelling from userdefined")
}

func (c *CompA) GetLiveData(index int) int {
	liveVersion := c.GetLiveVersion()

	switch liveVersion {
	case 0:
		return c.arrBlue[c.permBlue[index]]
	case 1:
		return c.arrGreen[c.permGreen[index]]
	default:
		return 0
	}

}

func (c *CompA) GetStagingData() interface{} {
	stagingVersion := c.GetStagingVersion()
	switch stagingVersion {
	case 0:
		return c.arrBlue
	case 1:
		return c.arrGreen
	default:
		return nil
	}

}

func (c *CompA) getReferences() *CompA {
	return c
}

func (c *CompA) GetStagingDatas() []int {
	stagingVersion := c.GetStagingVersion()
	switch stagingVersion {
	case 0:
		return c.permBlue
	case 1:
		return c.permGreen
	default:
		return nil
	}

}
