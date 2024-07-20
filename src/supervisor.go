package bg

// Supervisor represents the supervisor component that controls other components.
type Supervisor struct {
	Channel chan string
}

// NewSupervisor creates a new supervisor with a channel.
func NewSupervisor() *Supervisor {
	return &Supervisor{
		Channel: make(chan string),
	}
}
