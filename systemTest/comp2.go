package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type CompB struct {
	bg.BasicComponent
	arrBlue   []int
	arrGreen  []int
	permBlue  []int
	permGreen []int
	comp1Ref  *CompA // Reference to Comp1
}

// Function to calculate the permutation of an array
func (c *CompB) calculatePerm(arr []int) []int {
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

// Function to sum two arrays element by element
func (c *CompB) sumArrays(arr1, arr2 []int) []int {
	// Make sure both arrays are of the same length
	if len(arr1) != len(arr2) {
		fmt.Println("Arrays are not of the same length!")
		return nil
	}

	result := make([]int, len(arr1))
	for i := range arr1 {
		result[i] = arr1[i] + arr2[i]
	}
	return result
}

func (c *CompB) ProcessReq(request bg.Request[interface{}]) {
	//time.Sleep(300 * time.Millisecond)
	comp1Array := c.comp1Ref.GetStagingData().([]int)
	stagingVersion := c.GetStagingVersion()

	if request.ComponentName == reflect.TypeOf(c).Elem().Name() {
		if request.Operation == bg.Update {
			switch stagingVersion {
			case 0:
				c.arrBlue[request.Index] = request.Data.(int)
			case 1:
				c.arrGreen[request.Index] = request.Data.(int)
			default:
				fmt.Println("Unknown staging version:", stagingVersion)

			}

		}
	}

	var tempArray []int
	switch stagingVersion {
	case 0:
		tempArray = c.sumArrays(comp1Array, c.arrBlue)
		c.permBlue = c.calculatePerm(tempArray)
		//fmt.Println("Permutation of CompB: ", c.permBlue)
	case 1:
		tempArray = c.sumArrays(comp1Array, c.arrGreen)
		c.permGreen = c.calculatePerm(tempArray)
		//fmt.Println("Permutation of CompB: ", c.permGreen)
	default:
		fmt.Println("Unknown staging version:", stagingVersion)

	}

	fmt.Println("Process completed in CompB")

}

func (c *CompB) Sync() {
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
	fmt.Println("Sync completed in CompB")

}

func (c *CompB) Cancel() {
	fmt.Println("Component B cancelling")
	stagingVersion := c.GetStagingVersion()
	switch stagingVersion {
	case 0:
		c.arrBlue = c.arrGreen
		c.permBlue = c.permGreen
	case 1:
		c.arrGreen = c.arrBlue
		c.permGreen = c.permBlue
	default:
	}
}

func (c *CompB) GetLiveData(index int) int {
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

func (c *CompB) GetStagingData() interface{} {
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

func (c *CompB) getReferences() []bg.Component {
	var depCompRef []bg.Component
	depCompRef = append(depCompRef, c.comp1Ref.getReferences())
	depCompRef = append(depCompRef, c)

	return depCompRef

}

func (c *CompB) GetStagingDatas() []int {
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
