package bg

import (
	"context"
	"sync"
)

// Dependency represents a single dependency relationship between components.
type Dependency struct {
	Child  string
	Parent string
}

// InitializeComponents initializes and starts the supervisor and components based on dependencies.
func InitializeComponents(ctx context.Context, supervisor *Supervisor, names []string, dependencies []Dependency) map[string]Component {
	var wg sync.WaitGroup
	components := make(map[string]Component)

	// Create components
	for _, name := range names {
		components[name] = &BasicComponent{
			Name:         name,
			InChannel:    []chan string{},
			OutChannel:   []chan string{},
			SuperChannel: supervisor.Channel,
		}
	}

	// Set up channels based on dependencies
	for _, dep := range dependencies {
		parent := components[dep.Parent].(*BasicComponent)
		child := components[dep.Child].(*BasicComponent)
		newChannel := make(chan string)
		parent.OutChannel = append(parent.OutChannel, newChannel)
		child.InChannel = append(child.InChannel, newChannel)
	}

	// Start all components with the context
	for _, component := range components {
		component.Run(ctx, &wg)
	}

	// Ensure all goroutines are cleaned up before exiting
	go func() {
		wg.Wait()
	}()

	return components
}
