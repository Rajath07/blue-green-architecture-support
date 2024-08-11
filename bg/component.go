package bg

import (
	"context"
	"fmt"
	"sync"
)

// Component interface defines the behavior that all components must implement.
type Component interface {
	init(compId int, inChannel chan string)
	initOutChan(outChannel []chan string)
	getInChan() chan string
	sendSignal()
	run(ctx context.Context, wg *sync.WaitGroup)
	ProcessReq(ctx context.Context)
	CancelReq(ctx context.Context)
	SyncReq(ctx context.Context)
}

// BasicComponent represents a single component with channels and implements the Component interface.
type BasicComponent struct {
	CompId     int
	InChannel  chan string
	OutChannel []chan string
	//SuperChannel chan string
}

func (c *BasicComponent) init(compId int, inChannel chan string) {
	c.CompId = compId
	c.InChannel = make(chan string)
}

func (c *BasicComponent) initOutChan(outChannel []chan string) {
	c.OutChannel = outChannel
}

func (c *BasicComponent) getInChan() chan string {
	return c.InChannel
}

// Run starts the component's main execution loop.
func (c *BasicComponent) run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Component %d stopped due to cancellation\n", c.CompId)
				return
			case msg := <-c.InChannel:
				fmt.Printf("Component %d received message: %s\n", c.CompId, msg)
				compNameStructMap[msg].ProcessReq(ctx)
				compNameStructMap[msg].sendSignal()

			}
		}
	}()
}

func (c *BasicComponent) sendSignal() {
	//fmt.Println(c.OutChannel)
	//c.OutChannel[1] <- "Comp2"
	// for _, outChan := range c.OutChannel {
	// 	outChan <- "I am done"
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

// SyncReq is a placeholder for the user-defined request update method.
func (c *BasicComponent) SyncReq(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Printf("Component %d stopping request update due to cancellation\n", c.CompId)
		return
	default:
		fmt.Printf("Component %d updating request\n", c.CompId)
		// Example: Actual update logic
	}
}
