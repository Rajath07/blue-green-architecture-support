package bg

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	Channels map[int]chan string
}

// NewSupervisor creates a new supervisor with a channel.
func NewSupervisor(componentIds []int) *Supervisor {
	channels := make(map[int]chan string)
	for _, compId := range componentIds {
		channels[compId] = make(chan string)
	}
	return &Supervisor{Channels: channels}
}

// GetChannel returns the channel for the specified component.
func (s *Supervisor) GetChannel(componentId int) chan string {
	return s.Channels[componentId]
}
