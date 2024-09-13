package bg

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"gonum.org/v1/gonum/graph/simple"
)

type CompositeKey struct {
	myId   int
	compId int
}

var waitingCount = make(map[CompositeKey]int)
var idStructMap = map[int]Component{}
var waitCountSupervisor = make(map[int64]int)

// InitializeComponents initializes and starts the components based on dependencies.
func InitializeComponents(ctx context.Context, filePath string, userComps []Component, switchCount int) *Supervisor {
	var wg sync.WaitGroup
	var structNames []string
	//var idStructMap = map[int]Component{}
	var idInChanMap = make(map[int]chan interface{})
	var compNameStructMap = map[string]Component{}
	var superInChan = make(chan interface{})

	// Parse the YAML file
	redGraph, dependencies, err := ParseYAML(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Reduced graph ", dependencies)
	CountPaths(redGraph)
	waitCountSupervisor = calculateReachableNodes(redGraph)
	//fmt.Println("Waiting count ", waitingCount)
	//fmt.Println("Waiting count for supervisor ", waitCountSupervisor)

	// Assign IDs for components and store the names of the user defined structs
	for _, comp := range userComps {
		structName := reflect.TypeOf(comp).Elem().Name()
		structNames = append(structNames, structName)
		inChan := make(chan interface{})
		idInChanMap[getComponentId(structName)] = inChan
		comp.init(getComponentId(structName), inChan, false)
		idStructMap[getComponentId(structName)] = comp
		compNameStructMap[structName] = comp
	}

	// Assign outChannels based on dependencies
	for parent, children := range dependencies {
		var childInChannels = []chan interface{}{}
		childInChannels = append(childInChannels, superInChan) // Connect the supervisor input channel

		if len(children) == 0 {
			idStructMap[int(parent)].initOutChan(childInChannels) // If the parent has no children, set the outChannel to nil
		} else {
			for _, child := range children {
				childInChannels = append(childInChannels, idStructMap[int(child)].getInChan())
			}
			idStructMap[int(parent)].initOutChan(childInChannels)
		}
	}
	fmt.Println(IdToComponent)

	//Initialize the supervisor
	supervisor := initSupervisor(superInChan, idStructMap, switchCount)
	supervisor.run(ctx, &wg)

	// Start all components with the context
	for _, component := range userComps {
		component.run(ctx, &wg)
	}

	// Ensure all goroutines are cleaned up before exiting
	go func() {
		wg.Wait()
	}()

	return supervisor
}

// CountPaths calculates the waiting count for each component based on its dependencies
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

// DFSWithMemoization performs a DFS to calculate paths from the node to its ancestors.
func DFSWithMemoization(g *simple.DirectedGraph, nodeID int64, memo map[int64]map[int64]int) {
	// Check if the current node's ancestors are already calculated
	if _, ok := memo[nodeID]; ok {
		return
	}

	// Initialize the memo entry for the current node
	memo[nodeID] = make(map[int64]int)

	// Iterate over all predecessors (ancestors) of the current node
	to := g.To(nodeID)
	for to.Next() {
		pred := to.Node()
		predID := pred.ID()

		// Recursive DFS call for the predecessor
		DFSWithMemoization(g, predID, memo)

		// Update the count for each ancestor of the current node
		for ancestor, count := range memo[predID] {
			if count >= 1 {
				memo[nodeID][ancestor]++

			}
			// memo[nodeID][ancestor] += count
		}

		// Count the direct path from the predecessor to the current node
		memo[nodeID][predID]++
	}
}

// calculateReachableNodes calculates the number of reachable nodes for each node using memoization.
func calculateReachableNodes(graph *simple.DirectedGraph) map[int64]int {
	reachableCount := make(map[int64]int)
	memo := make(map[int64]bool)

	// Iterate through all nodes in the graph.
	nodes := graph.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		visited := make(map[int64]bool)
		// Compute the number of reachable nodes starting from this node.
		reachableCount[node.ID()] = dfsMemoized(graph, node.ID(), visited, memo)
	}

	return reachableCount
}

// dfsMemoized performs a DFS with memoization to count reachable nodes.
func dfsMemoized(graph *simple.DirectedGraph, nodeID int64, visited map[int64]bool, memo map[int64]bool) int {
	// If this node has already been visited, return 0 to avoid double-counting
	if visited[nodeID] {
		return 0
	}
	// Mark this node as visited
	visited[nodeID] = true

	// Initialize count with 1 to include this node itself.
	count := 1

	// Iterate over all successors (outgoing edges) of the current node.
	successors := graph.From(nodeID)
	for successors.Next() {
		neighbor := successors.Node()
		count += dfsMemoized(graph, neighbor.ID(), visited, memo)
	}

	return count
}
