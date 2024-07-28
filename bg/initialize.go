package bg

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// Dependency represents a single dependency relationship between components.
type Dependency struct {
	Child  int
	Parent int
}

type CompositeKey struct {
	myId   int
	compId int
}

var waitingCount = make(map[CompositeKey]int)
var graphNodes []graph.Node

// InitializeComponents initializes and starts the components based on dependencies.
func InitializeComponents(ctx context.Context, filePath string, userComps []Component) map[string]Component {
	var wg sync.WaitGroup
	var structNames []string
	var idStructMap = map[int]Component{}
	var idInChanMap = make(map[int]chan string)
	var compNameStructMap = map[string]Component{}

	// Parse the YAML file
	redGraph, dependencies, err := ParseYAML(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Reduced graph ", dependencies)
	CountPaths(redGraph)
	fmt.Println("Waiting count ", waitingCount)
	// Print the results
	for key, count := range waitingCount {
		fmt.Printf("Paths from node %d to node %d: %d\n", key.compId, key.myId, count)
	}

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
		component.run(ctx, &wg)
	}

	// Ensure all goroutines are cleaned up before exiting
	go func() {
		wg.Wait()
	}()

	return compNameStructMap
}

// CountPaths calculates the number of paths from each node to its ancestors
func CountPaths(g *simple.DirectedGraph) {
	memo := make(map[int64]map[int64]int)

	// Initialize the memoization map and perform DFS from each node
	nodes := g.Nodes()
	for nodes.Next() {
		nodeID := nodes.Node().ID()
		if _, ok := memo[nodeID]; !ok {
			DFSWithMemoization(g, nodeID, memo)
		}
	}

	// Populate the waitingCount map with the results from the memoization map
	for node, ancestors := range memo {
		for ancestor, count := range ancestors {
			key := CompositeKey{myId: int(node), compId: int(ancestor)}
			waitingCount[key] = count
		}
	}
}

// DFSWithMemoization performs a DFS and counts paths using memoization
func DFSWithMemoization(g *simple.DirectedGraph, nodeID int64, memo map[int64]map[int64]int) {
	// Check if the current node's ancestors are already calculated
	if _, ok := memo[nodeID]; ok {
		return
	}

	// Initialize the memo entry for the current node
	memo[nodeID] = make(map[int64]int)

	// Get the current node
	node := g.Node(nodeID)

	// Iterate over all predecessors (ancestors) of the current node
	to := g.To(node.ID())
	for to.Next() {
		pred := to.Node()
		predID := pred.ID()

		// Recursive DFS call for the predecessor
		DFSWithMemoization(g, predID, memo)

		// Update the count for each ancestor of the current node
		for ancestor, count := range memo[predID] {
			memo[nodeID][ancestor] += count
		}
		// Count the direct path from the predecessor to the current node
		memo[nodeID][predID]++
	}
}
