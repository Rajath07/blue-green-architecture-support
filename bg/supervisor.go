package bg

import (
	"context"
	"fmt"
	"sync"
)

type SupervisorInterface interface {
	run(ctx context.Context, wg *sync.WaitGroup)
	SendReq(componentName string, operation OperationType, data interface{})
}

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	CompId        int
	InChannel     chan string
	OutChannelMap map[int]chan string
}

// OperationType represents the type of operation for CRUD actions.
type OperationType int

const (
	Create OperationType = iota
	Read
	Update
	Delete
)

// Request encapsulates the details of a request being sent.
type Request[T any] struct {
	ComponentName string
	Operation     OperationType
	Data          T
}

// NewSupervisor creates a new supervisor with a channel.
func initSupervisor(inChan chan string, idStructMap map[int]Component) *Supervisor {
	var outChanMap = make(map[int]chan string)
	for id, comp := range idStructMap {
		outChanMap[id] = comp.getInChan()
	}
	return &Supervisor{CompId: 0, InChannel: inChan, OutChannelMap: outChanMap}
}

func (s *Supervisor) run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Component %d stopped due to cancellation\n", s.CompId)
				return
			case msg := <-s.InChannel:
				fmt.Printf("Supervisor received message: %s\n", msg)
			}
		}
	}()
}

// SendReq sends a request to a component, specifying the operation and data.
func (s *Supervisor) SendReq(componentName string, operation OperationType, data interface{}) {
	s.OutChannelMap[getComponentId(componentName)] <- componentName

}
