package bg

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// Component interface defines the behavior that all components must implement.
type Component interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	ProcessReq(ctx context.Context)
	CancelReq(ctx context.Context)
	SyncReq(ctx context.Context)
}

// BasicComponent represents a single component with channels and implements the Component interface.
type BasicComponent struct {
	CompId       int
	InChannel    []chan string
	OutChannel   []chan string
	SuperChannel chan string
}

// Run starts the component's main execution loop.
func (c *BasicComponent) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			selectCases := make([]reflect.SelectCase, len(c.InChannel)+2)
			for i, ch := range c.InChannel {
				selectCases[i] = reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(ch),
				}
			}
			// Add a case for ctx.Done() to stop the component gracefully
			selectCases[len(c.InChannel)] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ctx.Done()),
			}

			// Add a case for c.SuperChannel to handle supervisor signals
			selectCases[len(c.InChannel)+1] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c.SuperChannel),
			}

			chosen, value, ok := reflect.Select(selectCases)
			switch chosen {
			case len(c.InChannel): // ctx.Done() case
				fmt.Printf("Component %d stopped due to cancellation\n", c.CompId)
				//return
			case len(c.InChannel) + 1: // c.SuperChannel case
				if !ok {
					fmt.Printf("Component %d received from closed SuperChannel\n", c.CompId)
					continue
				}
				msg := value.String()
				fmt.Printf("Component %d received signal from supervisor: %s\n", c.CompId, msg)
				// Example: Handle supervisor signal
				// c.ProcessReq(ctx)
			default:
				if !ok {
					fmt.Printf("Component %d received from closed channel\n", c.CompId)
					continue
				}
				msg := value.String()
				fmt.Printf("Component %d received message: %s\n", c.CompId, msg)
			}
		}
	}()
}

// ProcessReq processes requests, checking for context cancellation.
func (c *BasicComponent) ProcessReq(ctx context.Context) {
	c.OutChannel[0] <- "Start Processing"
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

// UpdateReq is a placeholder for the user-defined request update method.
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
