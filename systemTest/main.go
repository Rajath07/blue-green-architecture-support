package main

import (
	"fmt"
	"time"

	"github.com/Rajath07/blue-green-architecture-support/bg"
)

func main() {
	compCollec := []bg.Component{}
	customComp1 := &CompA{arrBlue: []int{10, 45, 3, 76, 1}, arrGreen: []int{10, 45, 3, 76, 1}, permBlue: []int{4, 2, 0, 1, 3}, permGreen: []int{4, 2, 0, 1, 3}}
	customComp2 := &CompB{arrBlue: []int{8, 120, 32, 5, 1}, arrGreen: []int{8, 120, 32, 5, 1}, permBlue: []int{4, 0, 2, 3, 1}, permGreen: []int{4, 0, 2, 3, 1}, comp1Ref: customComp1}
	customComp3 := &CompC{arrBlue: []int{9, 99, 2, 3, 1}, arrGreen: []int{9, 99, 2, 3, 1}, permBlue: []int{4, 2, 0, 3, 1}, permGreen: []int{4, 2, 0, 3, 1}, comp1Ref: customComp1}
	customComp4 := &CompD{arrBlue: []int{0, 6, 4, 12, 5}, arrGreen: []int{0, 6, 4, 12, 5}, permBlue: []int{4, 0, 2, 3, 1}, permGreen: []int{4, 0, 2, 3, 1}, comp2Ref: customComp2}
	customComp5 := &CompE{arrBlue: []int{5, 2, 7, 3, 0}, arrGreen: []int{5, 2, 7, 3, 0}, permBlue: []int{4, 3, 0, 2, 1}, permGreen: []int{4, 3, 0, 2, 1}, comp2Ref: customComp2, comp3Ref: customComp3}
	customComp6 := &CompF{arrBlue: []int{1, 2, 3, 0, 400}, arrGreen: []int{1, 2, 3, 0, 400}, permBlue: []int{3, 0, 2, 1, 4}, permGreen: []int{3, 0, 2, 1, 4}, comp4Ref: customComp4, comp5Ref: customComp5}

	// {18, 165, 35, 81, 2}
	// {4, 0, 2, 3, 1}

	// {508, 165, 35, 81, 2}
	// {4, 2, 3, 1, 0}

	// comp4:
	// 	{
	// 	}
	// customComp5 := &Comp5{}
	compCollec = append(compCollec, customComp1)
	compCollec = append(compCollec, customComp2)
	compCollec = append(compCollec, customComp3)
	compCollec = append(compCollec, customComp4)
	compCollec = append(compCollec, customComp5)
	compCollec = append(compCollec, customComp6)
	// compCollec = append(compCollec, customComp5)
	components := bg.InitializeComponents("dependency.yaml", compCollec, 3)

	// fmt.Println("Initial permutation of CompB:", customComp2.GetStagingDatas())
	// fmt.Println("Initial permutation of CompD:", customComp4.GetStagingDatas())
	// fmt.Println("Initial permutation of CompE:", customComp5.GetStagingDatas())
	// fmt.Println("Initial permutation of CompF:", customComp6.GetStagingDatas())
	components.SendReq("CompB", bg.Update, 100, 0)
	components.SendReq("CompA", bg.Update, 500, 2)
	components.SendReq("CompF", bg.Update, 55, 3)
	time.Sleep(500 * time.Millisecond)

	// fmt.Println("Permutation of CompB before cancellation:", customComp2.GetStagingDatas())
	// fmt.Println("Permutation of CompD before cancellation:", customComp4.GetStagingDatas())
	// fmt.Println("Permutation of CompE before cancellation:", customComp5.GetStagingDatas())
	// fmt.Println("Permutation of CompF before cancellation:", customComp6.GetStagingDatas())
	//components.CancelReq("CompB")

	time.Sleep(1 * time.Second)
	// fmt.Println("Permutation of CompB:", customComp2.GetStagingDatas())
	// fmt.Println("Permutation of CompD:", customComp4.GetStagingDatas())
	// fmt.Println("Permutation of CompE:", customComp5.GetStagingDatas())
	// fmt.Println("Permutation of CompF:", customComp6.GetStagingDatas())
	// components.SendReq("CompB", bg.Update, 100, 0)
	// fmt.Println("Permutation of CompA:", customComp1.GetStagingDatas())

	// //time.Sleep(500 * time.Millisecond)
	// components.CancelReq("CompA")
	// time.Sleep(500 * time.Millisecond)
	// fmt.Println("Permutation of CompA:", customComp1.GetStagingDatas())

	fmt.Println("Reading CompA: ", customComp1.GetLiveData(4))
	fmt.Println("Reading CompB: ", customComp2.GetLiveData(4))
	fmt.Println("Reading CompC: ", customComp3.GetLiveData(4))
	fmt.Println("Reading CompD: ", customComp4.GetLiveData(4))
	fmt.Println("Reading CompE: ", customComp5.GetLiveData(4))
	fmt.Println("Reading CompF: ", customComp6.GetLiveData(4))

	// components.SendReq("CompD", bg.Update, 500, 0)

	time.Sleep(2 * time.Second)
	fmt.Println("Reading CompA: ", customComp1.GetLiveData(4))
	fmt.Println("Reading CompB: ", customComp2.GetLiveData(4))
	fmt.Println("Reading CompC: ", customComp3.GetLiveData(4))
	fmt.Println("Reading CompD: ", customComp4.GetLiveData(4))
	fmt.Println("Reading CompE: ", customComp5.GetLiveData(4))
	fmt.Println("Reading CompF: ", customComp6.GetLiveData(4))
	time.Sleep(20 * time.Second)

}
