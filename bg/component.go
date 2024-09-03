package bg

import (
	"context"
	"fmt"
	"sync"
)

// Component interface defines the behavior that all components must implement.
type Component interface {
	init(compId int, inChannel chan interface{}, DirtyFlag bool)
	initOutChan(outChannel []chan interface{})
	getInChan() chan interface{}
	sendSignal(req interface{}, state ComponentState, ctx context.Context)
	run(ctx context.Context, wg *sync.WaitGroup)
	getState() ComponentState
	setState(state ComponentState)
	GetLiveVersion() int
	GetStagingVersion() int
	ProcessReq(ctx context.Context)
	CancelReq(ctx context.Context)
	Switch(ctx context.Context)
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
	//SuperChannel chan string
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
func (c *BasicComponent) run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		var currCount int //To compare with the waiting count of the component
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Component %d stopped due to cancellation\n", c.CompId)
				return
			case msg := <-c.InChannel:
				switch request := msg.(type) {
				case Request[interface{}]:
					fmt.Printf("Component %d received request: %v\n", c.CompId, request)
					if request.ReqType == Operation {
						if request.SourceCompId == c.CompId {
							if component, exists := idStructMap[c.CompId]; exists {
								c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
								component.ProcessReq(ctx)
								c.DirtyFlag = true
								component.sendSignal(request, ComponentState(Idle), ctx)
								currCount = 0
							} else {
								fmt.Printf("Component %s not found in map\n", request.ComponentName)
							}
						} else {
							currCount++
							if currCount == waitingCount[CompositeKey{myId: c.CompId, compId: request.SourceCompId}] {
								if component, exists := idStructMap[c.CompId]; exists {
									//component.setState(Running)
									c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
									component.ProcessReq(ctx)
									c.DirtyFlag = true
									component.sendSignal(request, ComponentState(Idle), ctx)
									currCount = 0
								} else {
									fmt.Printf("Component %s not found in map\n", request.ComponentName)
								}
							} else {
								fmt.Println("Waiting to perform operation")
							}
						}
					} else if request.ReqType == Switch {
						if request.SourceCompId == c.CompId {
							if component, exists := idStructMap[c.CompId]; exists {
								c.OutChannel[0] <- Signal{SigType: request.ReqType, SourceCompId: request.SourceCompId, CompId: c.CompId, State: ComponentState(Running)}
								if c.DirtyFlag == true {
									component.Switch(ctx)
									c.DirtyFlag = false
									component.sendSignal(request, ComponentState(Idle), ctx)
									currCount = 0
									//Now send signal to others
								} else {
									fmt.Println("Component is not dirty")
									component.sendSignal(request, ComponentState(Idle), ctx)
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
										component.Switch(ctx)
										c.DirtyFlag = false
										component.sendSignal(request, ComponentState(Idle), ctx)
										currCount = 0
										//Now send signal to others
									} else {
										fmt.Println("Component is not dirty")
										component.sendSignal(request, ComponentState(Idle), ctx)
										currCount = 0
									}
								} else {
									fmt.Printf("Component not found in map\n")
								}
							} else {
								fmt.Println("Waiting to perform Switch")
							}
						}
					} else if request.ReqType == Cancel {

					}

				// case string:
				// 	if request == "Switch" {
				// 		if component, exists := idStructMap[c.CompId]; exists {
				// 			if c.DirtyFlag == true {
				// 				component.Switch(ctx)
				// 				c.DirtyFlag = false
				// 				//Now send signal to others
				// 			}

				// 		} else {
				// 			fmt.Printf("Component not found in map\n")
				// 		}
				// 	} else {
				// 		fmt.Printf("Component %d received an unexpected string: %s\n", c.CompId, request)
				// }

				default:
					fmt.Printf("Component %d received an unknown type of message\n", c.CompId)
				}
			default:
				if idStructMap[c.CompId].getState() == Cancelled {
					idStructMap[c.CompId].CancelReq(ctx)
					idStructMap[c.CompId].setState(Idle)
				}
			}
		}
	}()
}

func (c *BasicComponent) sendSignal(req interface{}, state ComponentState, ctx context.Context) {
	if c.getState() != Cancelled {
		switch v := req.(type) {
		case Request[interface{}]:
			c.OutChannel[0] <- Signal{SigType: v.ReqType, SourceCompId: v.SourceCompId, CompId: c.CompId, State: state}
			for i := 1; i < len(c.OutChannel); i++ { // Start from index 1
				c.OutChannel[i] <- req
			}
			// case string:
			// 	c.OutChannel[0] <- Signal{SourceCompId: -1, CompId: c.CompId, State: state}

		}

		// c.OutChannel[0] <- Signal{SourceCompId: req.SourceCompId, CompId: c.CompId, State: state}
		// for i := 1; i < len(c.OutChannel); i++ { // Start from index 1
		// 	c.OutChannel[i] <- req
		// }
	} else if c.getState() == Cancelled { //We have to call CancelReq here as well
		idStructMap[c.CompId].CancelReq(ctx)
		idStructMap[c.CompId].setState(Idle)
	}

}

// ProcessReq processes requests, checking for context cancellation.
func (c *BasicComponent) ProcessReq(ctx context.Context) {
	//c.OutChannel[1] <- "Start Processing"
	for _, outChan := range c.OutChannel {
		outChan <- "Start Processing"
	}
	select {
	case <-ctx.Done():
		fmt.Printf("Component %d stopping request processing due to cancellation\n", c.CompId)
		return
	default:
		fmt.Printf("Component %d processing request\n", c.CompId)
		// Example: Actual processing logic
	}
}

// CancelReq is a placeholder for the user-defined request cancellation method.
func (c *BasicComponent) CancelReq(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Printf("Component %d stopping request cancellation due to cancellation\n", c.CompId)
		return
	default:
		fmt.Printf("Component %d cancelling request\n", c.CompId)
		// Example: Actual cancellation logic
	}
}

// Switch is a placeholder for the user-defined request update method.
func (c *BasicComponent) Switch(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Printf("Component %d stopping request update due to cancellation\n", c.CompId)
		return
	default:
		fmt.Printf("Component %d updating request\n", c.CompId)
		// Example: Actual update logic
	}
}

// SetState sets the state of the component safely
func (c *BasicComponent) setState(state ComponentState) {
	c.StateMutex.Lock()
	defer c.StateMutex.Unlock()
	c.State = state
	//fmt.Println("Component", c.CompId, "state changed to ", state)
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
