package bg

import (
	"fmt"
	"os"

	"gonum.org/v1/gonum/graph/simple"
	"gopkg.in/yaml.v2"

	"gonum.org/v1/gonum/graph/traverse"

	"gonum.org/v1/gonum/graph"
)

type ComponentDependency struct {
	Components map[string][]string `yaml:"components"`
}

var IdToComponent = map[int]string{}
var ComponentToId = map[string]int{}

func getComponentId(name string) int {
	return ComponentToId[name]
}

func getComponentName(id int) string {
	return IdToComponent[id]
}

func createDependencyMap(g *simple.DirectedGraph) map[int64][]int64 {
	newMap := make(map[int64][]int64)

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
	edges := g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		from := edge.From()
		to := edge.To()

		// Temporarily remove the edge
		g.RemoveEdge(from.ID(), to.ID())

		// Check if there is still a path from 'from' to 'to'
		dfs := traverse.DepthFirst{}
		visited := make(map[int64]bool)
		dfs.Walk(g, from, func(n graph.Node) bool {
			visited[n.ID()] = true
			return false // continue walking
		})

		// If no path exists, re-add the edge
		if !visited[to.ID()] {
			g.SetEdge(edge)
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
	fmt.Println("\nTransitive Reduction Graph:")
	printGraph(graph)
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
