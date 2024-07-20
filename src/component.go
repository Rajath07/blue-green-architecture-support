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
	Name         string
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
			select {
			case <-ctx.Done():
				fmt.Printf("Component %s stopped due to cancellation\n", c.Name)
				return
			case msg := <-c.SuperChannel:
				fmt.Printf("Component %s received signal from supervisor: %s\n", c.Name, msg)
				// Example: Handle supervisor signal
				c.ProcessReq(ctx)
			default:
				// Listen on all InChannels using a dynamic select
				selectCases := make([]reflect.SelectCase, len(c.InChannel))
				for i, ch := range c.InChannel {
					selectCases[i] = reflect.SelectCase{
						Dir:  reflect.SelectRecv,
						Chan: reflect.ValueOf(ch),
					}
				}
				// Perform the select on all InChannels
				_, value, ok := reflect.Select(selectCases)
				if !ok {
					fmt.Printf("Component %s received from closed channel\n", c.Name)
					continue
				}
				msg := value.String()
				fmt.Printf("Component %s received message: %s\n", c.Name, msg)
			}
		}
	}()
}

// ProcessReq processes requests, checking for context cancellation.
func (c *BasicComponent) ProcessReq(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Printf("Component %s stopping request processing due to cancellation\n", c.Name)
		return
	default:
		fmt.Printf("Component %s processing request\n", c.Name)
		// Example: Actual processing logic
	}
}

// CancelReq is a placeholder for the user-defined request cancellation method.
func (c *BasicComponent) CancelReq(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Printf("Component %s stopping request cancellation due to cancellation\n", c.Name)
		return
	default:
		fmt.Printf("Component %s cancelling request\n", c.Name)
		// Example: Actual cancellation logic
	}
}

// UpdateReq is a placeholder for the user-defined request update method.
func (c *BasicComponent) SyncReq(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Printf("Component %s stopping request update due to cancellation\n", c.Name)
		return
	default:
		fmt.Printf("Component %s updating request\n", c.Name)
		// Example: Actual update logic
	}
}
