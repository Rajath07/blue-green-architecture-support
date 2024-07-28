package bg

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// Dependency represents a single dependency relationship between components.
type Dependency struct {
	Child  int
	Parent int
}

// InitializeComponents initializes and starts the components based on dependencies.
func InitializeComponents(ctx context.Context, filePath string, userComps []Component) map[string]Component {
	var wg sync.WaitGroup
	var structNames []string
	var idStructMap = map[int]Component{}
	var idInChanMap = make(map[int]chan string)
	var compNameStructMap = map[string]Component{}

	// Parse the YAML file
	dependencies, err := ParseYAML(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Reduced graph ", dependencies)

	// Assign IDs for components and store the names of the user defined structs
	for _, comp := range userComps {
		structName := reflect.TypeOf(comp).Elem().Name()
		structNames = append(structNames, structName)
		inChan := make(chan string)
		idInChanMap[getComponentId(structName)] = inChan
		comp.init(getComponentId(structName), inChan)
		idStructMap[getComponentId(structName)] = comp
		compNameStructMap[structName] = comp
	}

	// Assign outChannels based on dependencies
	for parent, children := range dependencies {
		var childInChannels = []chan string{}
		if len(children) == 0 {
			idStructMap[int(parent)].initOutChan(nil)
		} else {
			for _, child := range children {
				childInChannels = append(childInChannels, idStructMap[int(child)].getInChan())
			}
			idStructMap[int(parent)].initOutChan(childInChannels)
		}
	}

	// Start all components with the context
	for _, component := range userComps {
		component.Run(ctx, &wg)
	}

	// Ensure all goroutines are cleaned up before exiting
	go func() {
		wg.Wait()
	}()

	return compNameStructMap
}
