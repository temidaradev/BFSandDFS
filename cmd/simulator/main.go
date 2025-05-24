package main

import (
	"log"

	"bfsdfs/internal/simulator"
	"bfsdfs/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Create a new simulator with 10 nodes
	sim := simulator.NewSimulator(10)

	// Create a new game with the simulator
	game := ui.NewGame(sim)

	// Configure and run the game with larger window size
	ebiten.SetWindowSize(800, 700)
	ebiten.SetWindowTitle("BFS and DFS Graph Simulator")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
