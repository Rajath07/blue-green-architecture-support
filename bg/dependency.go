package bg

import (
	"fmt"
	"os"

	"gonum.org/v1/gonum/graph/simple"
	"gopkg.in/yaml.v2"

	"gonum.org/v1/gonum/graph"
)

type ComponentDependency struct {
	Components map[string][]string `yaml:"components"`
}

var IdToComponent = map[int]string{} //Gives the component name (string)
var ComponentToId = map[string]int{} //Gives the component id (int) when the name is passed

func getComponentId(name string) int {
	return ComponentToId[name]
}

func getComponentName(id int) string {
	return IdToComponent[id]
}

func createDependencyMap(g *simple.DirectedGraph) map[int64][]int64 {
	newMap := make(map[int64][]int64)
	// Initialize the map with all nodes to ensure each node is represented
	nodes := g.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		newMap[node.ID()] = []int64{}
	}

	// Process the edges to fill in dependencies
	edges := g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		fromID := edge.From().ID()
		toID := edge.To().ID()
		newMap[fromID] = append(newMap[fromID], toID)
	}

	return newMap
}

func transitiveReduction(g *simple.DirectedGraph) {
	// Step 1: Compute the reachability matrix using Floyd-Warshall algorithm
	n := g.Nodes().Len()
	reach := make([][]bool, n)
	for i := range reach {
		reach[i] = make([]bool, n)
	}

	nodes := g.Nodes()
	nodeIndex := make(map[int64]int)
	indexNode := make([]graph.Node, n)
	i := 0
	for nodes.Next() {
		node := nodes.Node()
		nodeIndex[node.ID()] = i
		indexNode[i] = node
		i++
	}

	edges := g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		from := nodeIndex[edge.From().ID()]
		to := nodeIndex[edge.To().ID()]
		reach[from][to] = true
	}

	// Floyd-Warshall algorithm to find reachability
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				reach[i][j] = reach[i][j] || (reach[i][k] && reach[k][j])
			}
		}
	}

	// Step 2: Identify and remove transitive edges
	edges = g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		from := nodeIndex[edge.From().ID()]
		to := nodeIndex[edge.To().ID()]

		// Check if there is an indirect path from 'from' to 'to'
		transitive := false
		for k := 0; k < n; k++ {
			if k != from && k != to && reach[from][k] && reach[k][to] {
				transitive = true
				break
			}
		}

		// If there is a transitive path, remove the edge
		if transitive {
			g.RemoveEdge(edge.From().ID(), edge.To().ID())
		}
	}
}

func ParseYAML(filePath string) (*simple.DirectedGraph, map[int64][]int64, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal YAML data into a Go struct
	var dependencies ComponentDependency
	err = yaml.Unmarshal(data, &dependencies)
	if err != nil {
		panic(err)
	}
	//comp1Slice := dependencies.Components["component1"]
	//fmt.Printf("Unmarshalled data: %+v\n", dependencies.Components)

	// Define IDs for components

	ComponentDependencyID := map[int][]int{}
	idCounter := 1
	for key, _ := range dependencies.Components {
		IdToComponent[idCounter] = key
		ComponentToId[key] = idCounter
		idCounter++
	}

	for key, child := range dependencies.Components {
		if len(child) == 0 {
			//ComponentDependencyID[getComponentId(key)] = []int{}
		} else {
			for _, childComp := range child {
				ComponentDependencyID[getComponentId(childComp)] = append(ComponentDependencyID[getComponentId(childComp)], getComponentId(key))
			}
		}
	}
	graph := simple.NewDirectedGraph()

	for key, _ := range ComponentDependencyID {
		graph.AddNode(simple.Node(key))
	}
	for key, _ := range ComponentDependencyID {
		for _, v := range ComponentDependencyID[key] {
			//graph.SetEdge(graph.NewEdge{F: simple.Node(key), T: simple.Node(v)})
			graph.SetEdge(graph.NewEdge(simple.Node(key), simple.Node(v)))
		}
	}

	//printGraph(graph)
	transitiveReduction(graph)
	//fmt.Println("\nTransitive Reduction Graph:")
	//printGraph(graph)
	// Create new dependency map from the reduced graph
	reducedDependencyMap := createDependencyMap(graph)
	//fmt.Println("\nNew Dependency Map: ", reducedDependencyMap)
	return graph, reducedDependencyMap, nil

}

func printGraph(g *simple.DirectedGraph) {
	fmt.Println("Nodes:")
	nodes := g.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		fmt.Printf("Node ID: %d\n", node.ID())
	}

	fmt.Println("\nEdges:")
	edges := g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		fmt.Printf("Edge from %d to %d\n", edge.From().ID(), edge.To().ID())
	}
}
