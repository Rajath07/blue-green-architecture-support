package bg

import (
	"context"
	"fmt"
	"sync"
)

type SupervisorInterface interface {
	run(ctx context.Context, wg *sync.WaitGroup)
	sendToAll()
}

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	CompId        int
	InChannel     chan string
	OutChannelMap map[int]chan string
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

func (s *Supervisor) SendToAll() {
	for _, outChan := range s.OutChannelMap {
		outChan <- "hello from supervisor"
	}
}
