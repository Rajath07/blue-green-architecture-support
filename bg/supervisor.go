package bg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SupervisorInterface interface {
	run(ctx context.Context, wg *sync.WaitGroup)
	SendReq(componentName string, operation OperationType, data interface{})
	processQueue()
}

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	CompId        int
	InChannel     chan interface{}
	OutChannelMap map[int]chan interface{}
	RequestQueue  []Request[interface{}]
	QueueMutex    sync.Mutex // To protect access to the queue
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
	SourceCompId  int
	ComponentName string
	Operation     OperationType
	Data          T
	Index         int
}

// NewSupervisor creates a new supervisor with a channel.
func initSupervisor(inChan chan interface{}, idStructMap map[int]Component) *Supervisor {
	var outChanMap = make(map[int]chan interface{})
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
				switch m := msg.(type) {
				case Request[interface{}]:
					fmt.Printf("Supervisor received request: %v\n", m)
				case Signal:
					fmt.Printf("Supervisor received signal: %v\n", m)
					component := idStructMap[m.SourceCompId]
					component.setState(m.State)
				}
				//fmt.Printf("Supervisor received message: %s\n", msg)
			default:
				s.processQueue()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// SendReq sends a request to a component, specifying the operation and data.
func (s *Supervisor) SendReq(componentName string, operation OperationType, data interface{}, index int) bool {
	componentId := getComponentId(componentName)
	req := Request[interface{}]{
		SourceCompId:  componentId,
		ComponentName: componentName,
		Operation:     operation,
		Data:          data,
		Index:         index,
	}
	// Lock the queue for thread-safe access
	s.QueueMutex.Lock()
	defer s.QueueMutex.Unlock()

	// Enqueue the request
	s.RequestQueue = append(s.RequestQueue, req)
	fmt.Printf("Enqueued request for component %s \n", componentName)
	return true

	// if outChan, ok := s.OutChannelMap[componentId]; ok {
	// 	outChan <- req
	// 	fmt.Printf("Sent request to component %s with operation %d and data %v, index %d\n", componentName, operation, data, index)
	// } else {
	// 	fmt.Printf("Component %s not found\n", componentName)
	// }
}

func (s *Supervisor) processQueue() {
	// Lock the queue to safely dequeue requests
	s.QueueMutex.Lock()
	defer s.QueueMutex.Unlock()

	if len(s.RequestQueue) > 0 {
		// Dequeue the first request
		req := s.RequestQueue[0]
		component := idStructMap[req.SourceCompId]
		// Check if the component is idle
		if component.getState() == Idle {
			// Dequeue the request if the component is idle
			s.RequestQueue = s.RequestQueue[1:]
			//set the component to running
			component.setState(Running)
			// Send the request to the appropriate component
			if outChan, ok := s.OutChannelMap[req.SourceCompId]; ok {
				outChan <- req
				fmt.Printf("Dispatched request to component %s\n", req.ComponentName)
			} else {
				fmt.Printf("Component %s not found\n", req.ComponentName)
			}
		}

	}
}
