package bg

import (
	"context"
	"sync"
	// "os"
	// "gopkg.in/yaml.v2"
)

// Dependency represents a single dependency relationship between components.
type Dependency struct {
	Child  int
	Parent int
}

// // ComponentConfig represents a component configuration.
// type ComponentConfig struct {
// 	Name string `yaml:"name"`
// }

// // Configuration represents the overall YAML configuration.
// type Configuration struct {
// 	Components   []ComponentConfig `yaml:"components"`
// 	Dependencies []Dependency      `yaml:"dependencies"`
// }

// InitializeComponents initializes and starts the components based on dependencies.
func InitializeComponents(ctx context.Context, compIds []int, dependencies []Dependency, userComps Component) map[int]Component {
	var wg sync.WaitGroup
	components := make(map[int]Component)

	// Create components
	for _, compId := range compIds {
		components[compId] = &BasicComponent{
			CompId:     compId,
			InChannel:  make(chan string),
			OutChannel: []chan string{},
			//SuperChannel: supervisor.GetChannel(compId),
		}
	}

	// Set up channels based on dependencies
	for _, dep := range dependencies {
		parent := components[dep.Parent].(*BasicComponent)
		child := components[dep.Child].(*BasicComponent)
		parent.OutChannel = append(parent.OutChannel, child.InChannel)
	}
	// testOutChannel := []chan string{}
	// testOutChannel = append(testOutChannel, components[1].(*BasicComponent).InChannel)
	// comp.init(10, make(chan string), testOutChannel)
	// comp.Run(ctx, &wg)
	// components[10] = comp

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

// LoadConfig loads the configuration from a YAML file
// func LoadConfig(configFile string) (*Configuration, error) {
// 	file, err := os.Open(configFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var config Configuration
// 	decoder := yaml.NewDecoder(file)
// 	err = decoder.Decode(&config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &config, nil
// }
