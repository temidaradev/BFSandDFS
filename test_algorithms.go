package main

import (
	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/graph"
	"fmt"
)

func main() {
	// Create a simple test graph
	g := graph.NewRandomGraph(5)

	fmt.Println("Testing Graph Algorithms")
	fmt.Println("========================")

	// Test Dijkstra's algorithm
	fmt.Println("\n1. Dijkstra's Algorithm:")
	neighbors := g.GetWeightedNeighbors()
	distances, predecessors := algorithms.Dijkstra(neighbors, 0, len(g.Nodes))
	for i, dist := range distances {
		fmt.Printf("  Distance to node %d: %.2f\n", i, dist)
	}
	fmt.Printf("  Predecessors: %v\n", predecessors)

	// Test A* algorithm
	fmt.Println("\n2. A* Algorithm:")
	positions := g.GetPositions()
	path, cost := algorithms.AStar(neighbors, 0, len(g.Nodes)-1, positions)
	fmt.Printf("  Path from 0 to %d: %v\n", len(g.Nodes)-1, path)
	fmt.Printf("  Total cost: %.2f\n", cost)

	// Test Topological Sort
	fmt.Println("\n3. Topological Sort:")
	unweightedNeighbors := g.GetUnweightedNeighbors()
	topOrder := algorithms.TopologicalSort(unweightedNeighbors, len(g.Nodes))
	fmt.Printf("  Topological order: %v\n", topOrder)

	// Test Kruskal's MST
	fmt.Println("\n4. Kruskal's MST:")
	mst := algorithms.Kruskal(g.WeightedEdges, len(g.Nodes))
	fmt.Printf("  MST edges: %d\n", len(mst))
	for _, edge := range mst {
		fmt.Printf("    %d -> %d (weight: %.2f)\n", edge.From, edge.To, edge.Weight)
	}

	// Test Prim's MST
	fmt.Println("\n5. Prim's MST:")
	primMST := algorithms.Prim(neighbors, len(g.Nodes))
	fmt.Printf("  MST edges: %d\n", len(primMST))
	for _, edge := range primMST {
		fmt.Printf("    %d -> %d (weight: %.2f)\n", edge.From, edge.To, edge.Weight)
	}

	// Test Tarjan's SCC
	fmt.Println("\n6. Tarjan's SCC:")
	sccs := algorithms.Tarjan(unweightedNeighbors, len(g.Nodes))
	fmt.Printf("  Found %d strongly connected components:\n", len(sccs))
	for i, scc := range sccs {
		fmt.Printf("    SCC %d: %v\n", i+1, scc)
	}

	// Test Kosaraju's SCC
	fmt.Println("\n7. Kosaraju's SCC:")
	kosarajuSCCs := algorithms.Kosaraju(unweightedNeighbors, len(g.Nodes))
	fmt.Printf("  Found %d strongly connected components:\n", len(kosarajuSCCs))
	for i, scc := range kosarajuSCCs {
		fmt.Printf("    SCC %d: %v\n", i+1, scc)
	}

	fmt.Println("\nAll algorithms tested successfully!")
}
