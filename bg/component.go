package bg

import (
	"fmt"
	"sync"
)

// Component interface defines the behavior that all components must implement.
type Component interface {
	init(compId int, inChannel chan interface{}, DirtyFlag bool)
	initOutChan(outChannel []chan interface{})
	getInChan() chan interface{}
	sendSignal(req interface{}, state ComponentState)
	run(wg *sync.WaitGroup)
	getState() ComponentState
	setState(state ComponentState)
	GetLiveVersion() int
	GetStagingVersion() int
	GetStagingData() interface{}
	ProcessReq(req Request[interface{}])
	Cancel()
	Sync()
}

// Signal represents a signal sent between component and Supervisor
type Signal struct {
	SigType      RequestType
	SourceCompId int
	CompId       int
	State        ComponentState
}

// ComponentState defines the possible states of a component
type ComponentState int

const (
	Idle ComponentState = iota
	Running
	Cancelled
)

// BasicComponent represents a single component with channels and implements the Component interface.
type BasicComponent struct {
	CompId     int
	InChannel  chan interface{}
	OutChannel []chan interface{}
	State      ComponentState // Field to track component state
	StateMutex sync.Mutex     // Mutex to protect state changes
	DirtyFlag  bool
}

func (c *BasicComponent) init(compId int, inChannel chan interface{}, DirtyFlag bool) {
	c.CompId = compId
	c.InChannel = make(chan interface{})
	c.State = Idle
	c.DirtyFlag = DirtyFlag
}

func (c *BasicComponent) initOutChan(outChannel []chan interface{}) {
	c.OutChannel = outChannel
}

func (c *BasicComponent) getInChan() chan interface{} {
	return c.InChannel
}

// Run starts the component's main execution loop.
func (c *BasicComponent) run(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		var currCount int //To compare with the waiting count of the component
		defer wg.Done()
		for {
			select {
			case msg := <-c.InChannel:
				switch request := msg.(type) {
				case Request[interface{}]:
					if request.ReqType == Operation {
						if request.SourceCompId == c.CompId {
							if component, exists := idStructMap[c.CompId]; exists {
								c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
								component.ProcessReq(request)
								c.DirtyFlag = true
								request.Data = nil //Remove the data after the source component is done processing so that the remaining components know that they were not the source component
								component.sendSignal(request, ComponentState(Idle))
								currCount = 0
							} else {
								fmt.Printf("Component %s not found in map\n", request.ComponentName)
							}
						} else {
							currCount++
							if currCount == waitingCount[CompositeKey{myId: c.CompId, compId: request.SourceCompId}] {
								if component, exists := idStructMap[c.CompId]; exists {
									c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
									component.ProcessReq(request)
									c.DirtyFlag = true
									component.sendSignal(request, ComponentState(Idle))
									currCount = 0
								} else {
									fmt.Printf("Component %s not found in map\n", request.ComponentName)
								}
							} else {
								//Waiting to perform operation
							}
						}
					} else if request.ReqType == Switch {
						if request.SourceCompId == c.CompId {
							if component, exists := idStructMap[c.CompId]; exists {
								c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
								if c.DirtyFlag == true {
									component.Sync()
									c.DirtyFlag = false
									component.sendSignal(request, ComponentState(Idle))
									currCount = 0
								} else {
									fmt.Println("Component ", getComponentName(c.CompId), "is not dirty")
									component.sendSignal(request, ComponentState(Idle))
									currCount = 0
								}
							} else {
								fmt.Printf("Component not found in map\n")
							}
						} else {
							currCount++
							if currCount == waitingCount[CompositeKey{myId: c.CompId, compId: request.SourceCompId}] {
								if component, exists := idStructMap[c.CompId]; exists {
									c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
									if c.DirtyFlag == true {
										component.Sync()
										c.DirtyFlag = false
										component.sendSignal(request, ComponentState(Idle))
										currCount = 0
									} else {
										fmt.Println("Component ", getComponentName(c.CompId), "is not dirty")
										component.sendSignal(request, ComponentState(Idle))
										currCount = 0
									}
								} else {
									fmt.Printf("Component not found in map\n")
								}
							} else {
								//Waiting to perform Switch
							}
						}
					}

				default:
					fmt.Printf("Component %d received an unknown type of message\n", c.CompId)
				}
			default:
				if idStructMap[c.CompId].getState() == Cancelled {
					idStructMap[c.CompId].Cancel()
					idStructMap[c.CompId].setState(Idle)
				}
			}
		}
	}()
}

func (c *BasicComponent) sendSignal(req interface{}, state ComponentState) {
	if c.getState() != Cancelled {
		switch v := req.(type) {
		case Request[interface{}]:
			c.OutChannel[0] <- Signal{SigType: v.ReqType, SourceCompId: v.SourceCompId, CompId: c.CompId, State: state}
			for i := 1; i < len(c.OutChannel); i++ { // Start from index 1
				c.OutChannel[i] <- req
			}
		}
	} else if c.getState() == Cancelled { //We have to call Cancel here as well
		idStructMap[c.CompId].Cancel()
		idStructMap[c.CompId].setState(Idle)
	}

}

// ProcessReq processes requests
func (c *BasicComponent) ProcessReq(req Request[interface{}]) {
	fmt.Printf("Component %d processing request\n", c.CompId)
	// Example: Actual processing logic
}

// Cancel is a placeholder for the user-defined request cancellation method.
func (c *BasicComponent) Cancel() {
	fmt.Printf("Component %d cancelling request\n", c.CompId)
	// Example: Actual cancellation logic
}

// Sync is a placeholder for the user-defined synchronization method.
func (c *BasicComponent) Sync() {
	fmt.Printf("Component %d updating request\n", c.CompId)
	// Example: Actual update logic
}

// SetState sets the state of the component safely
func (c *BasicComponent) setState(state ComponentState) {
	c.StateMutex.Lock()
	defer c.StateMutex.Unlock()
	c.State = state
}

// GetState gets the current state of the component safely
func (c *BasicComponent) getState() ComponentState {
	c.StateMutex.Lock()
	defer c.StateMutex.Unlock()
	return c.State
}

func (c *BasicComponent) GetLiveVersion() int {
	return int(liveVersion)
}

func (c *BasicComponent) GetStagingVersion() int {
	var liveVersion = c.GetLiveVersion()
	if liveVersion == int(Blue) {
		return int(Green)
	} else if liveVersion == int(Green) {
		return int(Blue)
	} else {
		return -1
	}
}

func (c *BasicComponent) GetStagingData() interface{} {
	return nil
}
