package bg

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	Channels map[string]chan string
}

// NewSupervisor creates a new supervisor with a channel.
func NewSupervisor(componentNames []string) *Supervisor {
	channels := make(map[string]chan string)
	for _, name := range componentNames {
		channels[name] = make(chan string)
	}
	return &Supervisor{Channels: channels}
}

// GetChannel returns the channel for the specified component.
func (s *Supervisor) GetChannel(componentName string) chan string {
	return s.Channels[componentName]
}
