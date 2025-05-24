package ui

import (
	"fmt"
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

	// Button background with rounded corners effect
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

	// Fill main background
	bg.Fill(buttonColor)

	// Add a subtle border
	borderColor := color.RGBA{
		uint8(max(int(buttonColor.R)-30, 0)),
		uint8(max(int(buttonColor.G)-30, 0)),
		uint8(max(int(buttonColor.B)-30, 0)),
		255,
	}

	// Draw border
	for i := 0; i < b.Width; i++ {
		bg.Set(i, 0, borderColor)          // Top
		bg.Set(i, b.Height-1, borderColor) // Bottom
	}
	for i := 0; i < b.Height; i++ {
		bg.Set(0, i, borderColor)         // Left
		bg.Set(b.Width-1, i, borderColor) // Right
	}

	// Draw button background
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(btnX), float64(btnY))
	screen.DrawImage(bg, opts)

	// Draw button text (centered)
	textWidth := len(b.Text) * 7 // Approximate width based on basicfont
	textX := btnX + (b.Width-textWidth)/2
	textY := btnY + b.Height/2 + 5 // +5 for centering with basicfont

	// Add text shadow for better visibility
	shadowColor := color.RGBA{0, 0, 0, 100}
	text.Draw(screen, b.Text, basicfont.Face7x13, textX+1, textY+1, shadowColor)
	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, b.TextColor)
}

// Game represents the Ebiten game that visualizes the simulation
type Game struct {
	Sim               *simulator.Simulator
	StartNode         int
	MouseX            int
	MouseY            int
	lastMouseX        int // Track last mouse X position for optimization
	lastMouseY        int // Track last mouse Y position for optimization
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

	// Performance optimization: cached images
	graphCanvas       *ebiten.Image
	gridCanvas        *ebiten.Image
	canvasNeedsRedraw bool
	lastGraphState    string  // Simple hash to track if graph state has changed
	lastCanvasOffsetX float64 // Track last canvas offset for optimization
	lastCanvasOffsetY float64 // Track last canvas offset for optimization

	// UI element caches
	textBgCache       *ebiten.Image
	sliderBgCache     *ebiten.Image
	sliderHandleCache *ebiten.Image
	speedBgCache      *ebiten.Image
	messageBgCache    *ebiten.Image

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

	ShowHelp bool // Add help overlay toggle

	// AVL Input Modal
	ShowAVLInput  bool
	AVLInputValue int
	AVLAction     string // "insert", "delete", "search"
	AVLInputText  string // Text input for AVL value

	// Selection features
	Selecting           bool
	SelectionStartX     int      // X position where selection drag started
	SelectionStartY     int      // Y position where selection drag started
	SelectedNodes       []int    // Indices of selected nodes
	SelectedEdges       [][2]int // Indices of selected edges (as pairs of node indices)
	DraggingSelection   bool     // Whether a selected group is being dragged
	SelectionDragStartX float64  // X position where dragging of selection started (canvas coords)
	SelectionDragStartY float64  // Y position where dragging of selection started (canvas coords)

	// Performance optimization fields
	lastFrameTime time.Time
	frameCount    int
	fps           int
	lastFPSUpdate time.Time
}

// NewGame creates a new game with the given simulator
func NewGame(sim *simulator.Simulator) *Game {
	// Get initial window size for canvas initialization
	screenWidth, screenHeight := ebiten.WindowSize()

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
		ShowHelp:       false, // Initialize help overlay as hidden

		// Initialize cached canvases
		graphCanvas:       ebiten.NewImage(screenWidth, screenHeight),
		gridCanvas:        ebiten.NewImage(screenWidth, screenHeight),
		canvasNeedsRedraw: true, // Force initial draw

		// Initialize UI element caches
		textBgCache:       ebiten.NewImage(screenWidth, 20),
		sliderBgCache:     ebiten.NewImage(200, 20),
		sliderHandleCache: ebiten.NewImage(10, 20),
		speedBgCache:      ebiten.NewImage(50, 20),
		messageBgCache:    ebiten.NewImage(200, 20),
	}

	// Create UI buttons
	g.createButtons()

	return g
}

// generateGraphStateHash creates a simple "hash" to detect graph state changes
func (g *Game) generateGraphStateHash() string {
	// This is a simple fingerprint of the current graph state
	// If this changes, we need to redraw the graph
	h := fmt.Sprintf("n%d-e%d-c%d-v%d-o%f-%f-g%v",
		len(g.Sim.Graph.Nodes),
		len(g.Sim.Graph.Edges),
		g.Sim.Current,
		len(g.Sim.Visited),
		g.CanvasOffsetX,
		g.CanvasOffsetY,
		g.ShowGrid)

	return h
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
	bottomRowY := 50 // BFS, DFS, AVL Tree, Step, Auto, Reset
	// algorithmRowY := 90 // Dijkstra, A*, Topo Sort, Kruskal, Prim, Tarjan, Kosaraju
	middleRowY := 90 // New Graph, Load, Save, Add Edge, Del Edge, Add Node, Del Node
	topRowY := 130   // Reset View, Grid, Snap, Edit Mode
	avlRowY := 170   // Insert, Delete, Search (AVL operations)

	// Create bottom row buttons - algorithm execution controls
	buttons := []*Button{
		{
			X: margin, Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "BFS", BgColor: blueBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle {
					if g.StartNode >= 0 && g.StartNode < len(g.Sim.Graph.Nodes) {
						g.Sim.StartBFS(g.StartNode)
						g.showMessage("BFS started from node " + string(rune('A'+g.StartNode)))
					} else {
						g.showMessage("Please select a start node first")
					}
				}
			},
		},
		{
			X: margin + (buttonWidth + buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "DFS", BgColor: blueBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle {
					if g.StartNode >= 0 && g.StartNode < len(g.Sim.Graph.Nodes) {
						g.Sim.StartDFS(g.StartNode)
						g.showMessage("DFS started from node " + string(rune('A'+g.StartNode)))
					} else {
						g.showMessage("Please select a start node first")
					}
				}
			},
		},
		{
			X: margin + 2*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "AVL Tree", BgColor: purpleBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeIdle {
					g.Sim.StartAVL()
					g.AutoStep = false // Disable auto-stepping in AVL mode
					g.showMessage("AVL Tree mode started. Use Insert/Delete/Search buttons.")
				} else if g.Sim.Mode == algorithms.ModeAVL {
					g.showMessage("Already in AVL Tree mode. Use Insert/Delete/Search buttons.")
				} else {
					g.showMessage("Reset first to switch to AVL Tree mode.")
				}
			},
		},
		{
			X: margin + 3*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Step", BgColor: greenBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Done {
					g.showMessage("Algorithm has completed. Reset to start over.")
				} else if g.Sim.Mode == algorithms.ModeIdle {
					g.showMessage("Please select an algorithm first.")
				} else if g.Sim.Mode == algorithms.ModeAVL {
					g.showMessage("Step not applicable in AVL Tree mode.")
				} else {
					g.Sim.Update()
					if g.Sim.Done {
						g.showMessage("Algorithm completed!")
					}
				}
			},
		},
		{
			X: margin + 4*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Auto", BgColor: orangeBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Done {
					g.showMessage("Algorithm has completed. Reset to start over.")
				} else if g.Sim.Mode == algorithms.ModeIdle {
					g.showMessage("Please select an algorithm first.")
				} else if g.Sim.Mode == algorithms.ModeAVL {
					g.showMessage("Auto stepping not applicable in AVL Tree mode.")
				} else {
					g.AutoStep = !g.AutoStep
					if g.AutoStep {
						g.showMessage("Auto stepping enabled. Use speed slider to adjust.")
					} else {
						g.showMessage("Auto stepping disabled.")
					}
				}
			},
		},
		{
			X: margin + 5*(buttonWidth+buttonSpacing), Y: bottomRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Reset", BgColor: redBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				g.Sim.Reset()
				g.AutoStep = false
				g.showMessage("Algorithm reset. Ready for new simulation.")
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
					// Create an empty graph instead of a random one
					g.Sim.Graph = graph.Graph{}
					g.Sim.Reset()
					g.StartNode = -1 // No start node for an empty graph initially
					g.showMessage("New empty graph created. Add nodes to start.")
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

	// Create AVL operation buttons - only shown when in AVL mode
	avlRowButtons := []*Button{
		{
			X: margin, Y: avlRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Insert", BgColor: greenBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeAVL {
					g.AVLAction = "insert"
					g.ShowAVLInput = true
					g.AVLInputText = ""
				}
			},
		},
		{
			X: margin + (buttonWidth + buttonSpacing), Y: avlRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Delete", BgColor: redBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeAVL {
					g.AVLAction = "delete"
					g.ShowAVLInput = true
					g.AVLInputText = ""
				}
			},
		},
		{
			X: margin + 2*(buttonWidth+buttonSpacing), Y: avlRowY, Width: buttonWidth, Height: buttonHeight,
			Text: "Search", BgColor: orangeBg, TextColor: whiteTxt, AnchorBottom: true,
			Action: func() {
				if g.Sim.Mode == algorithms.ModeAVL {
					g.AVLAction = "search"
					g.ShowAVLInput = true
					g.AVLInputText = ""
				}
			},
		},
	}

	// Add buttons to the game
	//buttons = append(buttons, algorithmRowButtons...)
	buttons = append(buttons, middleRowButtons...)
	buttons = append(buttons, topRowButtons...)
	buttons = append(buttons, avlRowButtons...)
	g.Buttons = buttons
}

// showMessage displays a temporary message to the user
func (g *Game) showMessage(msg string) {
	g.Message = msg
	g.MessageTimer = 120 // Display for 2 seconds (120 frames at 60 FPS)
}

// Layout returns the game's logical screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Call resize handler
	g.HandleResize(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
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

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
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

	// Mark canvas for redraw
	g.canvasNeedsRedraw = true
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

	// Mark canvas for redraw
	g.canvasNeedsRedraw = true
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

	// Mark canvas for redraw
	g.canvasNeedsRedraw = true
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

		// Mark canvas for redraw
		g.canvasNeedsRedraw = true

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

	// Mark canvas for redraw
	g.canvasNeedsRedraw = true
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
	if ebiten.IsKeyPressed(ebiten.KeyA) && !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle && g.Sim.Mode != algorithms.ModeAVL {
		g.AutoStep = !g.AutoStep
		// Wait to avoid repeated toggles
		time.Sleep(200 * time.Millisecond)
	}

	// Step key (space)
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle && g.Sim.Mode != algorithms.ModeAVL {
		g.Sim.Update()
		// Wait to avoid too-rapid stepping
		time.Sleep(100 * time.Millisecond)
	}
}
