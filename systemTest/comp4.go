package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

type CompD struct {
	bg.BasicComponent
	arrBlue   []int
	arrGreen  []int
	permBlue  []int
	permGreen []int
	comp2Ref  *CompB // Reference to Comp1
}

// Function to calculate the permutation of an array
func (c *CompD) calculatePerm(arr []int) []int {
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
func (c *CompD) sumArrays(arr1, arr2 []int) []int {
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

func (c *CompD) ProcessReq(request bg.Request[interface{}]) {
	//time.Sleep(1 * time.Second)
	depCompRefs := c.comp2Ref.getReferences()
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

	// Retrieve arrays from all ancestor components
	var totalLength int
	ancestorArrays := [][]int{}

	for _, depComp := range depCompRefs {
		ancestorArray := depComp.GetStagingData().([]int) // Call GetStagingData on each ancestor
		if ancestorArray != nil {
			ancestorArrays = append(ancestorArrays, ancestorArray)
			totalLength += len(ancestorArray)
		}
	}

	var ownArray []int
	switch stagingVersion {
	case 0:
		ownArray = c.arrBlue
	case 1:
		ownArray = c.arrGreen
	default:
		fmt.Println("Unknown staging version:", stagingVersion)
		return
	}

	// Create a temporary array that sums corresponding elements from ancestor arrays and ownArray
	tempSumArray := make([]int, len(ownArray))
	for i := range ownArray {
		sum := ownArray[i]
		// Sum corresponding elements from each ancestor's array
		for _, ancestorArray := range ancestorArrays {
			if i < len(ancestorArray) {
				sum += ancestorArray[i]
			}
		}
		tempSumArray[i] = sum
	}

	switch stagingVersion {
	case 0:
		c.permBlue = c.calculatePerm(tempSumArray)
	case 1:
		c.permGreen = c.calculatePerm(tempSumArray)
	default:
		fmt.Println("Unknown staging version:", stagingVersion)

	}
	fmt.Println("Process completed in CompD")
}

func (c *CompD) Sync() {
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
	fmt.Println("Sync completed in CompD")

}

func (c *CompD) Cancel() {
	fmt.Println("Component D cancelling")
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

func (c *CompD) GetLiveData(index int) int {
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

func (c *CompD) GetStagingData() interface{} {
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

func (c *CompD) getReferences() []bg.Component {
	var depCompRef []bg.Component
	depCompRef = append(depCompRef, c.comp2Ref.getReferences()...)
	depCompRef = append(depCompRef, c)

	return depCompRef

}

func (c *CompD) GetStagingDatas() []int {
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
