package simulator

import (
	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/graph"
)

// Simulator manages the graph traversal simulation
type Simulator struct {
	Graph      graph.Graph
	Order      []int
	Queue      []int
	Stack      []int
	Visited    map[int]bool
	Current    int // Currently active node
	LastActive int // Last active node for animation
	Mode       algorithms.TraversalMode
	Step       int
	Done       bool
}

// NewSimulator creates a new simulator with n nodes
func NewSimulator(n int) *Simulator {
	return &Simulator{
		Graph:      graph.NewRandomGraph(n),
		Visited:    map[int]bool{},
		Current:    -1,
		LastActive: -1,
		Mode:       algorithms.ModeIdle,
	}
}

// Update performs one step of the selected algorithm
func (s *Simulator) Update() error {
	if s.Done || s.Mode == algorithms.ModeIdle {
		return nil
	}

	// Create a map of node neighbors for the algorithm functions
	neighbors := make(map[int][]int, len(s.Graph.Nodes))
	for i, node := range s.Graph.Nodes {
		neighbors[i] = node.Neighbors
	}

	var nextNode int
	var isDone bool

	// Update the last active node
	s.LastActive = s.Current

	if s.Mode == algorithms.ModeBFS {
		s.Queue, nextNode, isDone = algorithms.BFSStep(s.Queue, s.Visited, neighbors)
	} else if s.Mode == algorithms.ModeDFS {
		s.Stack, nextNode, isDone = algorithms.DFSStep(s.Stack, s.Visited, neighbors)
	}

	s.Done = isDone

	// If we found a new node to visit, add it to the order
	if nextNode != -1 {
		s.Current = nextNode
		s.Visited[nextNode] = true
		s.Order = append(s.Order, nextNode)
	}

	return nil
}

// StartBFS starts a BFS traversal from the given start node
func (s *Simulator) StartBFS(start int) {
	s.Mode = algorithms.ModeBFS
	s.Queue = []int{start}
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false
}

// StartDFS starts a DFS traversal from the given start node
func (s *Simulator) StartDFS(start int) {
	s.Mode = algorithms.ModeDFS
	s.Stack = []int{start}
	s.Queue = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false
}

// Reset clears the simulation state
func (s *Simulator) Reset() {
	s.Mode = algorithms.ModeIdle
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false
}
