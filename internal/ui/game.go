package ui

import (
	"image/color"
	"time"

	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/graph"
	"bfsdfs/internal/simulator"
	"bfsdfs/pkg/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Button represents a clickable UI button
type Button struct {
	X, Y, Width, Height int
	Text                string
	BgColor             color.RGBA
	TextColor           color.RGBA
	Hover               bool
	Action              func()
	// Anchor properties for HUD positioning
	AnchorRight  bool // If true, X position is calculated from right edge
	AnchorBottom bool // If true, Y position is calculated from bottom edge
}

// IsInside checks if coordinates are inside the button
func (b *Button) IsInside(x, y int) bool {
	// Note: This will be used less frequently after we implement the getAdjustedButtonPosition method
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

// Draw renders the button on the screen
func (b *Button) Draw(screen *ebiten.Image, g *Game) {
	var btnX, btnY int

	// Get adjusted position based on anchoring if Game is provided
	if g != nil {
		btnX, btnY = g.getAdjustedButtonPosition(b)
	} else {
		btnX, btnY = b.X, b.Y
	}

	// Button background
	bg := ebiten.NewImage(b.Width, b.Height)

	// Different color for hover state
	buttonColor := b.BgColor
	if b.Hover {
		// Lighten the color for hover effect
		buttonColor = color.RGBA{
			uint8(min(int(buttonColor.R)+40, 255)),
			uint8(min(int(buttonColor.G)+40, 255)),
			uint8(min(int(buttonColor.B)+40, 255)),
			buttonColor.A,
		}
	}

	bg.Fill(buttonColor)

	// Draw button background
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(btnX), float64(btnY))
	screen.DrawImage(bg, opts)

	// Draw button text (centered)
	textWidth := len(b.Text) * 7 // Approximate width based on basicfont
	textX := btnX + (b.Width-textWidth)/2
	textY := btnY + b.Height/2 + 5 // +5 for centering with basicfont
	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, b.TextColor)
}

// Game represents the Ebiten game that visualizes the simulation
type Game struct {
	Sim               *simulator.Simulator
	StartNode         int
	MouseX            int
	MouseY            int
	MouseClicked      bool
	MouseReleased     bool
	MouseRightClicked bool
	AutoStep          bool
	StepDelay         int  // Frames to wait between auto-steps
	StepCounter       int  // Current frame count for auto-stepping
	SliderDragging    bool // Whether the speed slider is being dragged
	Buttons           []*Button

	// Node editing features
	EditMode      bool
	DraggingNode  int // Index of node being dragged, -1 if none
	AddingEdge    bool
	EdgeStartNode int
	RemovingNode  bool
	RemovingEdge  bool

	// Grid features
	ShowGrid   bool
	SnapToGrid bool
	GridConfig draw.GridConfig

	// Canvas movement features
	CanvasOffsetX    float64 // X offset for canvas movement
	CanvasOffsetY    float64 // Y offset for canvas movement
	CanvasDragging   bool    // Whether the canvas is being dragged
	CanvasDragStartX int     // X position where canvas drag started
	CanvasDragStartY int     // Y position where canvas drag started

	// Context menu
	ContextMenu *ContextMenu

	// File dialogs
	SaveDialog     *FileDialog
	LoadDialog     *FileDialog
	ShowSaveDialog bool
	ShowLoadDialog bool

	// Message display
	Message      string
	MessageTimer int
}

// NewGame creates a new game with the given simulator
func NewGame(sim *simulator.Simulator) *Game {
	g := &Game{
		Sim:            sim,
		StartNode:      0,
		StepDelay:      30, // Default to 30 frames between steps (about 0.5 seconds at 60 FPS)
		DraggingNode:   -1, // No node being dragged initially
		EdgeStartNode:  -1, // No edge start node selected initially
		ShowGrid:       true,
		SnapToGrid:     true,
		GridConfig:     draw.DefaultGridConfig(),
		ContextMenu:    NewContextMenu(),
		SaveDialog:     NewFileDialog(true),
		LoadDialog:     NewFileDialog(false),
		CanvasOffsetX:  0, // Initial canvas offset
		CanvasOffsetY:  0, // Initial canvas offset
		CanvasDragging: false,
	}

	// Create UI buttons
	g.createButtons()

	return g
}

// createButtons initializes all UI buttons
func (g *Game) createButtons() {
	// Button colors
	blueBg := color.RGBA{70, 130, 180, 255}    // Steel blue
	greenBg := color.RGBA{60, 160, 60, 255}    // Green
	redBg := color.RGBA{180, 60, 60, 255}      // Red
	orangeBg := color.RGBA{220, 130, 30, 255}  // Orange
	purpleBg := color.RGBA{130, 60, 180, 255}  // Purple
	grayBg := color.RGBA{100, 100, 110, 255}   // Gray
	whiteTxt := color.RGBA{240, 240, 240, 255} // White text

	// Button dimensions
	buttonWidth := 80
	buttonHeight := 30
	buttonSpacing := 10
	margin := 20

	// Fixed positions for button rows (bottom to top)
	bottomRowY := 50 // BFS, DFS, Step, Auto, Reset
	middleRowY := 90 // New Graph, Load, Save, Add Edge, Del Edge, Add Node, Del Node
	topRowY := 130   // Reset View, Grid, Snap, Edit Mode

	// Create bottom row buttons - algorithm execution controls
	buttons := []*Button{
		{
			X: margin, Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "BFS", BgColor: blueBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle {
					g.Sim.StartBFS(g.StartNode)
				}
			},
		},
		{
			X: margin + (buttonWidth + buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "DFS", BgColor: blueBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle {
					g.Sim.StartDFS(g.StartNode)
				}
			},
		},
		{
			X: margin + 2*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Step", BgColor: greenBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle {
					g.Sim.Update()
				}
			},
		},
		{
			X: margin + 3*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Auto", BgColor: orangeBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle {
					g.AutoStep = !g.AutoStep
				}
			},
		},
		{
			X: margin + 4*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Reset", BgColor: redBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.Sim.Reset()
				g.AutoStep = false
			},
		},
	}

	// Create middle row buttons - graph modification controls
	middleRowButtons := []*Button{
		{
			X: margin, Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "New Graph", BgColor: purpleBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle {
					nodeCount := len(g.Sim.Graph.Nodes)
					*g.Sim = *simulator.NewSimulator(nodeCount)
					g.StartNode = 0
				}
			},
		},
		{
			X: margin + (buttonWidth + buttonSpacing), Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Load", BgColor: color.RGBA{160, 120, 60, 255}, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.LoadDialog.Show()
				g.ShowLoadDialog = true
			},
		},
		{
			X: margin + 2*(buttonWidth+buttonSpacing), Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Save", BgColor: color.RGBA{60, 120, 160, 255}, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.SaveDialog.Show()
				g.ShowSaveDialog = true
			},
		},
		{
			X: margin + 3*(buttonWidth+buttonSpacing), Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Add Edge", BgColor: blueBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.AddingEdge = !g.AddingEdge
				g.RemovingNode = false
				g.RemovingEdge = false
				g.EditMode = g.AddingEdge

				if g.AddingEdge {
					g.showMessage("Click two nodes to add an edge between them")
				}
			},
		},
		{
			X: margin + 4*(buttonWidth+buttonSpacing), Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Del Edge", BgColor: orangeBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.RemovingEdge = !g.RemovingEdge
				g.AddingEdge = false
				g.RemovingNode = false
				g.EditMode = g.RemovingEdge

				if g.RemovingEdge {
					g.showMessage("Click on two nodes to remove the edge between them")
				}
			},
		},
		{
			X: margin + 5*(buttonWidth+buttonSpacing), Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Add Node", BgColor: greenBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle && len(g.Sim.Graph.Nodes) < 15 {
					nodeCount := len(g.Sim.Graph.Nodes)
					*g.Sim = *simulator.NewSimulator(nodeCount + 1)
					g.StartNode = 0
				}
			},
		},
		{
			X: margin + 6*(buttonWidth+buttonSpacing), Y: middleRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Del Node", BgColor: redBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.RemovingNode = !g.RemovingNode
				g.AddingEdge = false
				g.RemovingEdge = false
				g.EditMode = g.RemovingNode

				if g.RemovingNode {
					g.showMessage("Click a node to remove it")
				}
			},
		},
	}

	// Create top row buttons - view controls
	topRowButtons := []*Button{
		{
			X: margin, Y: topRowY, Width: buttonWidth + 20, Height: buttonHeight,
			Text: "Reset View", BgColor: color.RGBA{80, 80, 140, 255}, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.CanvasOffsetX = 0
				g.CanvasOffsetY = 0
				g.showMessage("Canvas view reset")
			},
		},
		{
			X: margin + (buttonWidth + 20 + buttonSpacing), Y: topRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Grid", BgColor: grayBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.ShowGrid = !g.ShowGrid
				if g.ShowGrid {
					g.showMessage("Grid display enabled")
				} else {
					g.showMessage("Grid display disabled")
				}
			},
		},
		{
			X: margin + (buttonWidth + 20 + buttonSpacing) + (buttonWidth + buttonSpacing), Y: topRowY,
			Width: buttonWidth, Height: buttonHeight,
			Text: "Snap", BgColor: grayBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.SnapToGrid = !g.SnapToGrid
				if g.SnapToGrid {
					g.showMessage("Snap to grid enabled")
				} else {
					g.showMessage("Snap to grid disabled")
				}
			},
		},
		{
			X: margin + (buttonWidth + 20 + buttonSpacing) + 2*(buttonWidth+buttonSpacing), Y: topRowY,
			Width: buttonWidth + 20, Height: buttonHeight,
			Text: "Edit Mode", BgColor: grayBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.EditMode = !g.EditMode
				if !g.EditMode {
					g.RemovingNode = false
					g.AddingEdge = false
					g.RemovingEdge = false
				}

				if g.EditMode {
					g.showMessage("Edit mode: Drag nodes to reposition them")
				}
			},
		},
	}

	// Add buttons to the game
	buttons = append(buttons, middleRowButtons...)
	buttons = append(buttons, topRowButtons...)
	g.Buttons = buttons
}

// showMessage displays a temporary message to the user
func (g *Game) showMessage(msg string) {
	g.Message = msg
	g.MessageTimer = 120 // Display for 2 seconds (120 frames at 60 FPS)
}

// Note: The Update method implementation has been moved to updater.go
// to avoid duplicate method definitions

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper functions for graph editing
func (g *Game) addNode(x, y int) {
	// Create a new node
	newNode := graph.Node{
		X:         x,
		Y:         y,
		Neighbors: []int{},
	}

	// Add to the simulator's graph
	g.Sim.Graph.Nodes = append(g.Sim.Graph.Nodes, newNode)
}

func (g *Game) removeNode(index int) {
	// First remove any edges connected to this node
	newEdges := [][2]int{}
	for _, edge := range g.Sim.Graph.Edges {
		if edge[0] != index && edge[1] != index {
			// Adjust indices for nodes after the removed one
			e := [2]int{edge[0], edge[1]}
			if e[0] > index {
				e[0]--
			}
			if e[1] > index {
				e[1]--
			}
			newEdges = append(newEdges, e)
		}
	}
	g.Sim.Graph.Edges = newEdges

	// Update neighbors lists
	for i := range g.Sim.Graph.Nodes {
		if i == index {
			continue
		}

		// Remove the deleted node from neighbor lists
		newNeighbors := []int{}
		for _, nb := range g.Sim.Graph.Nodes[i].Neighbors {
			if nb != index {
				// Adjust indices for nodes after the removed one
				if nb > index {
					nb--
				}
				newNeighbors = append(newNeighbors, nb)
			}
		}
		g.Sim.Graph.Nodes[i].Neighbors = newNeighbors
	}

	// Remove the node itself
	g.Sim.Graph.Nodes = append(g.Sim.Graph.Nodes[:index], g.Sim.Graph.Nodes[index+1:]...)

	// Adjust start node if necessary
	if g.StartNode == index {
		g.StartNode = 0
	} else if g.StartNode > index {
		g.StartNode--
	}
}

func (g *Game) addEdge(a, b int) {
	// Check if the edge already exists
	for _, edge := range g.Sim.Graph.Edges {
		if (edge[0] == a && edge[1] == b) || (edge[0] == b && edge[1] == a) {
			return // Edge already exists
		}
	}

	// Add the new edge
	g.Sim.Graph.Edges = append(g.Sim.Graph.Edges, [2]int{a, b})

	// Update neighbors
	g.Sim.Graph.Nodes[a].Neighbors = append(g.Sim.Graph.Nodes[a].Neighbors, b)
	g.Sim.Graph.Nodes[b].Neighbors = append(g.Sim.Graph.Nodes[b].Neighbors, a)
}

func (g *Game) removeEdge(a, b int) {
	// Find and remove the edge
	edgeIndex := -1
	for i, edge := range g.Sim.Graph.Edges {
		if (edge[0] == a && edge[1] == b) || (edge[0] == b && edge[1] == a) {
			edgeIndex = i
			break
		}
	}

	if edgeIndex != -1 {
		g.Sim.Graph.Edges = append(g.Sim.Graph.Edges[:edgeIndex], g.Sim.Graph.Edges[edgeIndex+1:]...)

		// Update node neighbors
		g.removeFromNeighbors(a, b)
		g.removeFromNeighbors(b, a)

		g.showMessage("Edge removed")
	} else {
		g.showMessage("No edge exists between these nodes")
	}
}

func (g *Game) removeFromNeighbors(nodeIndex, neighborToRemove int) {
	neighbors := g.Sim.Graph.Nodes[nodeIndex].Neighbors
	newNeighbors := []int{}

	for _, n := range neighbors {
		if n != neighborToRemove {
			newNeighbors = append(newNeighbors, n)
		}
	}

	g.Sim.Graph.Nodes[nodeIndex].Neighbors = newNeighbors
}

// clearNodeEdges removes all edges connected to a specific node
func (g *Game) clearNodeEdges(nodeIndex int) {
	// Find and remove all edges connected to this node
	newEdges := [][2]int{}
	for _, edge := range g.Sim.Graph.Edges {
		if edge[0] != nodeIndex && edge[1] != nodeIndex {
			newEdges = append(newEdges, edge)
		}
	}
	g.Sim.Graph.Edges = newEdges

	// Remove the node from all other nodes' neighbor lists
	for i := range g.Sim.Graph.Nodes {
		if i == nodeIndex {
			// Clear this node's neighbors completely
			g.Sim.Graph.Nodes[i].Neighbors = []int{}
			continue
		}

		// Remove the node from this node's neighbors
		newNeighbors := []int{}
		for _, nb := range g.Sim.Graph.Nodes[i].Neighbors {
			if nb != nodeIndex {
				newNeighbors = append(newNeighbors, nb)
			}
		}
		g.Sim.Graph.Nodes[i].Neighbors = newNeighbors
	}
}

// handleKeyboardInput maintains keyboard control support for convenience
func handleKeyboardInput(g *Game) {
	// BFS key
	if ebiten.IsKeyPressed(ebiten.KeyB) && g.Sim.Mode == algorithms.ModeIdle {
		g.Sim.StartBFS(g.StartNode)
	}

	// DFS key
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.Sim.Mode == algorithms.ModeIdle {
		g.Sim.StartDFS(g.StartNode)
	}

	// Reset key
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Sim.Reset()
		g.AutoStep = false
	}

	// Toggle auto-step (A key)
	if ebiten.IsKeyPressed(ebiten.KeyA) && !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle {
		g.AutoStep = !g.AutoStep
		// Wait to avoid repeated toggles
		time.Sleep(200 * time.Millisecond)
	}

	// Step key (space)
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle {
		g.Sim.Update()
		// Wait to avoid too-rapid stepping
		time.Sleep(100 * time.Millisecond)
	}
}

// Layout returns the game's logical screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600 // Increased from 500, 600 to give more width for buttons
}

// getAdjustedButtonPosition calculates the button position based on anchoring
func (g *Game) getAdjustedButtonPosition(btn *Button) (int, int) {
	btnX := btn.X
	btnY := btn.Y

	// Adjust X position for right-anchored buttons
	if btn.AnchorRight {
		// Get screen width and adjust from right edge
		w, _ := ebiten.WindowSize()
		btnX = w - btn.X - btn.Width
	}

	// Adjust Y position for bottom-anchored buttons
	if btn.AnchorBottom {
		// Get screen height and adjust from bottom edge
		_, h := ebiten.WindowSize()
		btnY = h - btn.Y - btn.Height
	}

	return btnX, btnY
}
