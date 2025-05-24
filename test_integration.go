package main

import (
	"fmt"
	"log"

	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/simulator"
)

func main() {
	fmt.Println("Testing Integration with Simulator")
	fmt.Println("=================================")

	// Create a new simulator with 5 nodes
	sim := simulator.NewSimulator(5)

	// Test all algorithm start methods
	fmt.Println("\n1. Testing Start Methods...")

	// Test Dijkstra
	sim.StartDijkstra(0)
	if sim.Mode == algorithms.ModeDijkstra && sim.Done {
		fmt.Println("✓ StartDijkstra working")
	} else {
		log.Fatal("✗ StartDijkstra failed")
	}

	// Reset and test A*
	sim.Reset()
	sim.StartAStar(0, 4)
	if sim.Mode == algorithms.ModeAStar && sim.Done {
		fmt.Println("✓ StartAStar working")
	} else {
		log.Fatal("✗ StartAStar failed")
	}

	// Reset and test Topological
	sim.Reset()
	sim.StartTopological()
	if sim.Mode == algorithms.ModeTopological && sim.Done {
		fmt.Println("✓ StartTopological working")
	} else {
		log.Fatal("✗ StartTopological failed")
	}

	// Reset and test Kruskal
	sim.Reset()
	sim.StartKruskal()
	if sim.Mode == algorithms.ModeKruskal && sim.Done {
		fmt.Println("✓ StartKruskal working")
	} else {
		log.Fatal("✗ StartKruskal failed")
	}

	// Reset and test Prim
	sim.Reset()
	sim.StartPrim()
	if sim.Mode == algorithms.ModePrim && sim.Done {
		fmt.Println("✓ StartPrim working")
	} else {
		log.Fatal("✗ StartPrim failed")
	}

	// Reset and test Tarjan
	sim.Reset()
	sim.StartTarjan()
	if sim.Mode == algorithms.ModeTarjan && sim.Done {
		fmt.Println("✓ StartTarjan working")
	} else {
		log.Fatal("✗ StartTarjan failed")
	}

	// Reset and test Kosaraju
	sim.Reset()
	sim.StartKosaraju()
	if sim.Mode == algorithms.ModeKosaraju && sim.Done {
		fmt.Println("✓ StartKosaraju working")
	} else {
		log.Fatal("✗ StartKosaraju failed")
	}

	fmt.Println("\n2. Testing Result Getter Methods...")

	// Test result getters (after running Dijkstra)
	sim.StartDijkstra(0)
	distances := sim.GetShortestPaths()
	if distances != nil {
		fmt.Println("✓ GetShortestPaths working")
	} else {
		log.Fatal("✗ GetShortestPaths failed")
	}

	// Test path getter (after running A*)
	sim.StartAStar(0, 4)
	path := sim.GetPath()
	if path != nil {
		fmt.Println("✓ GetPath working")
	} else {
		log.Fatal("✗ GetPath failed")
	}

	// Test MST getter (after running Kruskal)
	sim.StartKruskal()
	mst := sim.GetMST()
	if mst != nil {
		fmt.Println("✓ GetMST working")
	} else {
		log.Fatal("✗ GetMST failed")
	}

	// Test SCC getter (after running Tarjan)
	sim.StartTarjan()
	sccs := sim.GetSCCs()
	if sccs != nil {
		fmt.Println("✓ GetSCCs working")
	} else {
		log.Fatal("✗ GetSCCs failed")
	}

	// Test topological order getter
	sim.StartTopological()
	topOrder := sim.GetTopologicalOrder()
	if topOrder != nil {
		fmt.Println("✓ GetTopologicalOrder working")
	} else {
		log.Fatal("✗ GetTopologicalOrder failed")
	}

	fmt.Println("\n✅ All integration tests passed!")
	fmt.Println("The advanced graph algorithms are now fully integrated with the simulator.")
	fmt.Println("You can run the application with: go run ./cmd/simulator")
}
