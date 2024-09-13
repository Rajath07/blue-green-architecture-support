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
	CancelReq(componentName string)
	processQueue()
	updateTaskList()
}

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	CompId        int
	InChannel     chan interface{}
	OutChannelMap map[int]chan interface{}
	RequestQueue  []Request[interface{}]
	QueueMutex    sync.Mutex    // To protect access to the queue
	TaskList      map[int][]int //Maintains current ongoing tasks
	DoneList      map[int][]int
	SwitchList    map[int][]int //Contains the list of components that are done switching
	switchCount   int
}

// OperationType represents the type of operation for CRUD actions.
type OperationType int

const (
	Create OperationType = iota
	Update
	Delete
)

type LiveVersion int

const (
	Blue LiveVersion = iota
	Green
)

type RequestType int

const (
	Operation RequestType = iota
	Switch
	Cancel
)

// Request encapsulates the details of a request being sent.
type Request[T any] struct {
	ReqType       RequestType
	SourceCompId  int
	ComponentName string
	Operation     OperationType
	Data          T
	Index         int
}

type CompRequest[T any] struct {
	ComponentName string
	Operation     OperationType
	Data          T
	Index         int
}

var switchCount = 0
var liveVersion = Blue
var versionToggled = false

// NewSupervisor creates a new supervisor with a channel.
func initSupervisor(inChan chan interface{}, idStructMap map[int]Component, switchCount int) *Supervisor {
	var outChanMap = make(map[int]chan interface{})
	for id, comp := range idStructMap {
		outChanMap[id] = comp.getInChan()
	}
	return &Supervisor{CompId: 0, InChannel: inChan, OutChannelMap: outChanMap, RequestQueue: []Request[interface{}]{}, TaskList: make(map[int][]int), DoneList: make(map[int][]int), SwitchList: make(map[int][]int), switchCount: switchCount}
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
					//fmt.Printf("Supervisor received signal: %v\n", m)
					component := idStructMap[m.SourceCompId]
					component.setState(m.State)
					s.updateTaskList(m)

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
		ReqType:       Operation,
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
	fmt.Printf("Enqueued request for component %d \n", componentId)
	return true
}

func (s *Supervisor) CancelReq(componentName string) {
	for _, compIds := range s.TaskList[getComponentId(componentName)] {
		idStructMap[compIds].setState(Cancelled)

	}
	delete(s.TaskList, getComponentId(componentName))
	delete(s.DoneList, getComponentId(componentName))
}

func (s *Supervisor) processQueue() {
	// Lock the queue to safely dequeue requests
	s.QueueMutex.Lock()
	defer s.QueueMutex.Unlock()

	if switchCount == s.switchCount {
		if versionToggled != true {
			//Toggle liveVersion when switchCount is reached
			if liveVersion == Blue {
				liveVersion = Green
			} else {
				liveVersion = Blue
			}
			versionToggled = true
		}
		// var component Component
		// var compId int
		//If doneList is empty then we reset the switchCount to 0
		if len(s.DoneList) == 0 && len(s.TaskList) == 0 {
			switchCount = 0
			versionToggled = false
		} else {
			if len(s.TaskList) == 0 {
				for sourceCompId, _ := range s.DoneList {
					req := Request[interface{}]{ReqType: Switch, SourceCompId: sourceCompId}
					s.TaskList[sourceCompId] = []int{}
					s.SwitchList[sourceCompId] = []int{}
					if outChan, ok := s.OutChannelMap[sourceCompId]; ok {
						outChan <- req
						fmt.Printf("Dispatched switch signal to component %d\n", sourceCompId)
					}
					delete(s.DoneList, sourceCompId)
					break
				}

			}
			// compToSwitch := s.DoneList[0]
			// s.DoneList = s.DoneList[1:]
			// fmt.Println("compToSwitch ", compToSwitch)
			//Extract key of this map
			// for key, v := range s.DoneList {
			// 	// compId := k
			// 	// component := idStructMap[k]
			// 	// fmt.Println("CompID to be switched ", k)
			// 	for _, compId := range v {
			// 		if outChan, ok := s.OutChannelMap[compId]; ok {
			// 			outChan <- "Switch"
			// 			fmt.Printf("Dispatched switch signal to component %d\n", compId)
			// 		}
			// 	}
			// 	delete(s.DoneList, key)
			// }
			// if component.getState() == Idle && s.TaskList[compId] == nil {
			// 	s.DoneList = s.DoneList[1:]
			// 	s.TaskList[compId] = []int{}
			// 	component.setState(Running)
			// 	//Now send the switch signal to the component
			// 	if outChan, ok := s.OutChannelMap[compId]; ok {
			// 		outChan <- Signal{SourceCompId: s.CompId, CompId: compId, State: Idle}
			// 		fmt.Printf("Dispatched switch signal to component %d\n", compId)
			// 	} else {
			// 		fmt.Printf("Component %d not found\n", compId)
			// 	}
			// }

		}

	} else if len(s.RequestQueue) > 0 {
		// Dequeue the first request
		req := s.RequestQueue[0]
		component := idStructMap[req.SourceCompId]
		// Check if the component is idle and there is no entry for that component in task list
		if component.getState() == Idle && len(s.TaskList) == 0 {
			s.RequestQueue = s.RequestQueue[1:] // Dequeue the request if the component is idle
			s.TaskList[req.SourceCompId] = []int{}
			s.DoneList[req.SourceCompId] = []int{}
			//component.setState(Running)

			// Send the request to the appropriate component
			if outChan, ok := s.OutChannelMap[req.SourceCompId]; ok {
				outChan <- req
				//fmt.Printf("Dispatched request to component %s\n", req.ComponentName)
			} else {
				fmt.Printf("Component %s not found\n", req.ComponentName)
			}
		}
		// else if s.TaskList[req.SourceCompId] != nil {
		// 	for _, compIds := range s.TaskList[req.SourceCompId] {
		// 		idStructMap[compIds].setState(Cancelled)
		// 	}
		// 	delete(s.TaskList, req.SourceCompId)
		// 	delete(s.DoneList, req.SourceCompId)

		// }

	}
}

func (s *Supervisor) updateTaskList(m Signal) {

	// if _, ok := s.TaskList[m.SourceCompId]; !ok {
	// 	s.TaskList[m.CompId] = []int{}
	// }
	if m.State == Running {
		s.TaskList[m.SourceCompId] = append(s.TaskList[m.SourceCompId], m.CompId)

	}
	if m.State == Idle && m.SigType == Operation {
		s.DoneList[m.SourceCompId] = append(s.DoneList[m.SourceCompId], m.CompId)
		if len(s.DoneList[m.SourceCompId]) == waitCountSupervisor[int64(m.SourceCompId)] {
			//fmt.Println("Deleting TaskList entry of Source Component ID", m.SourceCompId)
			//s.DoneList = append(s.DoneList, map[int][]int{m.SourceCompId: s.TaskList[m.SourceCompId]}) //Move the TaskList entry to the done list
			delete(s.TaskList, m.SourceCompId)
			switchCount++
		}
	}
	if m.State == Idle && m.SigType == Switch {
		s.SwitchList[m.SourceCompId] = append(s.SwitchList[m.SourceCompId], m.CompId)
		if len(s.SwitchList[m.SourceCompId]) == waitCountSupervisor[int64(m.SourceCompId)] {
			delete(s.SwitchList, m.SourceCompId)
			delete(s.TaskList, m.SourceCompId)
		}
	}
	// fmt.Println("TaskList", s.TaskList)
	// fmt.Println("SwitchList", s.SwitchList)
	// fmt.Println("DoneList", s.DoneList)

}
