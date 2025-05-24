package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"bfsdfs/internal/algorithms"
)

// Node represents a vertex in a graph with positioning information
type Node struct {
	X, Y      int
	Neighbors []int
	Weights   []float64 // Weights corresponding to neighbors
}

// Graph represents a collection of nodes and edges
type Graph struct {
	Nodes         []Node
	Edges         [][2]int
	WeightedEdges []algorithms.Edge // For weighted algorithms
}

// NewRandomGraph creates a new graph with n nodes and random edges
func NewRandomGraph(n int) Graph {
	g := Graph{}

	// Create nodes in a grid layout
	for i := 0; i < n; i++ {
		node := Node{
			X:         60 + (i%5)*80,
			Y:         60 + (i/5)*80,
			Neighbors: []int{},
			Weights:   []float64{},
		}
		g.Nodes = append(g.Nodes, node)
	}

	// Generate random edges
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	edgeSet := map[[2]int]bool{}

	// Create approximately n*2 edges
	for i := 0; i < n*2; i++ {
		a := r.Intn(n)
		b := r.Intn(n)

		// Avoid self-loops and duplicate edges
		if a != b && !edgeSet[[2]int{a, b}] && !edgeSet[[2]int{b, a}] {
			// Generate random weight between 1 and 10
			weight := 1.0 + r.Float64()*9.0

			// Add edges in both directions (undirected graph)
			g.Nodes[a].Neighbors = append(g.Nodes[a].Neighbors, b)
			g.Nodes[a].Weights = append(g.Nodes[a].Weights, weight)
			g.Nodes[b].Neighbors = append(g.Nodes[b].Neighbors, a)
			g.Nodes[b].Weights = append(g.Nodes[b].Weights, weight)

			// Store the edge for drawing
			g.Edges = append(g.Edges, [2]int{a, b})
			g.WeightedEdges = append(g.WeightedEdges, algorithms.Edge{From: a, To: b, Weight: weight})
			edgeSet[[2]int{a, b}] = true
		}
	}

	return g
}

// SaveGraph saves a graph to a JSON file
func (g *Graph) SaveGraph(filename string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Marshal the graph to JSON
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal graph: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadGraph loads a graph from a JSON file
func LoadGraph(filename string) (*Graph, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("graph file does not exist")
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON to graph
	var g Graph
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph: %w", err)
	}

	return &g, nil
}

// GetSavedGraphs returns a list of available saved graph filenames
func GetSavedGraphs(directory string) ([]string, error) {
	// Create directory if it doesn't exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		return []string{}, nil
	}

	// List all .json files in directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Filter for .json files
	var graphFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			graphFiles = append(graphFiles, filepath.Join(directory, file.Name()))
		}
	}

	return graphFiles, nil
}

// GetWeightedNeighbors returns a map of weighted neighbors for algorithms
func (g *Graph) GetWeightedNeighbors() map[int][]algorithms.Edge {
	neighbors := make(map[int][]algorithms.Edge)
	for i, node := range g.Nodes {
		neighbors[i] = []algorithms.Edge{}
		for j, neighbor := range node.Neighbors {
			weight := 1.0 // Default weight
			if j < len(node.Weights) {
				weight = node.Weights[j]
			}
			neighbors[i] = append(neighbors[i], algorithms.Edge{
				From:   i,
				To:     neighbor,
				Weight: weight,
			})
		}
	}
	return neighbors
}

// GetPositions returns a map of node positions for A* algorithm
func (g *Graph) GetPositions() map[int]algorithms.Position {
	positions := make(map[int]algorithms.Position)
	for i, node := range g.Nodes {
		positions[i] = algorithms.Position{X: node.X, Y: node.Y}
	}
	return positions
}

// GetUnweightedNeighbors returns a map of unweighted neighbors for algorithms
func (g *Graph) GetUnweightedNeighbors() map[int][]int {
	neighbors := make(map[int][]int)
	for i, node := range g.Nodes {
		neighbors[i] = node.Neighbors
	}
	return neighbors
}
