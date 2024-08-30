package bg

import (
	"context"
	"fmt"
	"sync"
)

// Component interface defines the behavior that all components must implement.
type Component interface {
	init(compId int, inChannel chan interface{})
	initOutChan(outChannel []chan interface{})
	getInChan() chan interface{}
	sendSignal(req interface{}, state ComponentState)
	run(ctx context.Context, wg *sync.WaitGroup)
	getState() ComponentState
	setState(state ComponentState)
	ProcessReq(ctx context.Context)
	CancelReq(ctx context.Context)
	Switch(ctx context.Context)
}

// Signal represents a signal sent between component and Supervisor
type Signal struct {
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
	//SuperChannel chan string
}

func (c *BasicComponent) init(compId int, inChannel chan interface{}) {
	c.CompId = compId
	c.InChannel = make(chan interface{})
	c.State = Idle
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
				if request, ok := msg.(Request[interface{}]); ok {
					fmt.Printf("Component %d received request: %v\n", c.CompId, request)
					if request.SourceCompId == c.CompId {
						if component, exists := idStructMap[c.CompId]; exists {
							component.ProcessReq(ctx)
							component.sendSignal(request, ComponentState(Idle))
							currCount = 0
						} else {
							fmt.Printf("Component %s not found in map\n", request.ComponentName)
						}

					} else {
						currCount++
						if currCount == waitingCount[CompositeKey{myId: c.CompId, compId: request.SourceCompId}] {
							// Check if the component exists in the map
							if component, exists := idStructMap[c.CompId]; exists {
								component.ProcessReq(ctx)
								component.sendSignal(request, ComponentState(Idle))
								currCount = 0
							} else {
								fmt.Printf("Component %s not found in map\n", request.ComponentName)
							}
						} else {
							fmt.Println("Waiting")
						}
					}

				} else {
					if component, exists := idStructMap[c.CompId]; exists {
						component.Switch(ctx)
					} else {
						fmt.Printf("Component %s not found in map\n", request.ComponentName)
					}
				}

			}
		}
	}()
}

func (c *BasicComponent) sendSignal(req interface{}, state ComponentState) {
	switch v := req.(type) {
	case Request[interface{}]:
		c.OutChannel[0] <- Signal{SourceCompId: v.SourceCompId, CompId: c.CompId, State: state}
		for i := 1; i < len(c.OutChannel); i++ { // Start from index 1
			c.OutChannel[i] <- req
		}
	case string:
		c.OutChannel[0] <- Signal{SourceCompId: -1, CompId: c.CompId, State: state}

	}

	// c.OutChannel[0] <- Signal{SourceCompId: req.SourceCompId, CompId: c.CompId, State: state}
	// for i := 1; i < len(c.OutChannel); i++ { // Start from index 1
	// 	c.OutChannel[i] <- req
	// }
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
