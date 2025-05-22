package simulator

import (
	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/graph"
)

// Simulator represents the graph traversal simulator
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
	avlTree    *algorithms.AVLTree
	avlValue   int
	avlAction  string // "insert", "delete", "search"
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

// StartAVL initializes the simulator for AVL tree operations
func (s *Simulator) StartAVL() {
	s.Mode = algorithms.ModeAVL
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false
	s.avlTree = algorithms.NewAVLTree()
	s.avlValue = 0
	s.avlAction = "insert"
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
	} else if s.Mode == algorithms.ModeAVL {
		// AVL specific update logic will go here
		// For now, we can just set Done to true to prevent infinite loops
		// or handle it based on AVL operation steps.
		isDone = true // Placeholder
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

// UpdateAVL updates the AVL tree visualization
func (s *Simulator) UpdateAVL() {
	if s.avlTree == nil {
		return
	}

	// Update node positions for visualization
	s.avlTree.UpdatePositions(400, 50, 60) // Center of screen, starting Y, level height
}

// InsertAVL inserts a value into the AVL tree
func (s *Simulator) InsertAVL(value int) {
	if s.avlTree == nil {
		s.StartAVL()
	}
	s.avlTree.Insert(value)
	s.UpdateAVL()
}

// DeleteAVL deletes a value from the AVL tree
func (s *Simulator) DeleteAVL(value int) {
	if s.avlTree == nil {
		return
	}
	s.avlTree.Delete(value)
	s.UpdateAVL()
}

// SearchAVL searches for a value in the AVL tree
func (s *Simulator) SearchAVL(value int) *algorithms.AVLNode {
	if s.avlTree == nil {
		return nil
	}
	return s.avlTree.Search(value)
}

// GetMode returns the current mode
func (s *Simulator) GetMode() algorithms.TraversalMode {
	return s.Mode
}

// GetAVLTree returns the AVL tree
func (s *Simulator) GetAVLTree() *algorithms.AVLTree {
	return s.avlTree
}

// GetAVLAction returns the current AVL action
func (s *Simulator) GetAVLAction() string {
	return s.avlAction
}

// SetAVLAction sets the current AVL action
func (s *Simulator) SetAVLAction(action string) {
	s.avlAction = action
}

// GetAVLValue returns the current AVL value
func (s *Simulator) GetAVLValue() int {
	return s.avlValue
}

// IncrementAVLValue increments the AVL value
func (s *Simulator) IncrementAVLValue() {
	s.avlValue++
}

// DecrementAVLValue decrements the AVL value
func (s *Simulator) DecrementAVLValue() {
	if s.avlValue > 0 {
		s.avlValue--
	}
}

// SetAVLValue sets the value for AVL operations
func (s *Simulator) SetAVLValue(value int) {
	s.avlValue = value
}
