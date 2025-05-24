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

	// Algorithm-specific results
	ShortestPaths map[int]float64
	Predecessors  map[int]int
	Path          []int
	MST           []algorithms.Edge
	SCCs          [][]int
	TopOrder      []int
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

// StartDijkstra initializes Dijkstra's algorithm from a source node
func (s *Simulator) StartDijkstra(source int) {
	s.Mode = algorithms.ModeDijkstra
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = source
	s.LastActive = -1
	s.Done = false

	// Run Dijkstra's algorithm
	neighbors := s.Graph.GetWeightedNeighbors()
	distances, predecessors := algorithms.Dijkstra(neighbors, source, len(s.Graph.Nodes))
	s.ShortestPaths = distances
	s.Predecessors = predecessors
	s.Done = true
}

// StartAStar initializes A* algorithm from source to goal
func (s *Simulator) StartAStar(source, goal int) {
	s.Mode = algorithms.ModeAStar
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = source
	s.LastActive = -1
	s.Done = false

	// Run A* algorithm
	neighbors := s.Graph.GetWeightedNeighbors()
	positions := s.Graph.GetPositions()
	path, _ := algorithms.AStar(neighbors, source, goal, positions)
	s.Path = path
	s.Done = true
}

// StartTopological initializes topological sort
func (s *Simulator) StartTopological() {
	s.Mode = algorithms.ModeTopological
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false

	// Run topological sort
	neighbors := s.Graph.GetUnweightedNeighbors()
	topOrder := algorithms.TopologicalSort(neighbors, len(s.Graph.Nodes))
	s.TopOrder = topOrder
	s.Done = true
}

// StartKruskal initializes Kruskal's MST algorithm
func (s *Simulator) StartKruskal() {
	s.Mode = algorithms.ModeKruskal
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false

	// Run Kruskal's algorithm
	mst := algorithms.Kruskal(s.Graph.WeightedEdges, len(s.Graph.Nodes))
	s.MST = mst
	s.Done = true
}

// StartPrim initializes Prim's MST algorithm
func (s *Simulator) StartPrim() {
	s.Mode = algorithms.ModePrim
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false

	// Run Prim's algorithm
	neighbors := s.Graph.GetWeightedNeighbors()
	mst := algorithms.Prim(neighbors, len(s.Graph.Nodes))
	s.MST = mst
	s.Done = true
}

// StartTarjan initializes Tarjan's SCC algorithm
func (s *Simulator) StartTarjan() {
	s.Mode = algorithms.ModeTarjan
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false

	// Run Tarjan's algorithm
	neighbors := s.Graph.GetUnweightedNeighbors()
	sccs := algorithms.Tarjan(neighbors, len(s.Graph.Nodes))
	s.SCCs = sccs
	s.Done = true
}

// StartKosaraju initializes Kosaraju's SCC algorithm
func (s *Simulator) StartKosaraju() {
	s.Mode = algorithms.ModeKosaraju
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false

	// Run Kosaraju's algorithm
	neighbors := s.Graph.GetUnweightedNeighbors()
	sccs := algorithms.Kosaraju(neighbors, len(s.Graph.Nodes))
	s.SCCs = sccs
	s.Done = true
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

	// Clear algorithm results
	s.ShortestPaths = nil
	s.Predecessors = nil
	s.Path = nil
	s.MST = nil
	s.SCCs = nil
	s.TopOrder = nil
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

// GetShortestPaths returns the shortest paths from Dijkstra
func (s *Simulator) GetShortestPaths() map[int]float64 {
	return s.ShortestPaths
}

// GetPath returns the path found by A*
func (s *Simulator) GetPath() []int {
	return s.Path
}

// GetMST returns the minimum spanning tree edges
func (s *Simulator) GetMST() []algorithms.Edge {
	return s.MST
}

// GetSCCs returns the strongly connected components
func (s *Simulator) GetSCCs() [][]int {
	return s.SCCs
}

// GetTopologicalOrder returns the topological ordering
func (s *Simulator) GetTopologicalOrder() []int {
	return s.TopOrder
}

// resetState resets common simulation state
func (s *Simulator) resetState() {
	s.Queue = nil
	s.Stack = nil
	s.Order = nil
	s.Visited = map[int]bool{}
	s.Current = -1
	s.LastActive = -1
	s.Done = false
	s.ShortestPaths = nil
	s.Predecessors = nil
	s.Path = nil
	s.MST = nil
	s.SCCs = nil
	s.TopOrder = nil
}
