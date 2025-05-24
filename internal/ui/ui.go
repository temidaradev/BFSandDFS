package ui

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"

	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/simulator"
)

// UI represents the user interface
type UI struct {
	simulator *simulator.Simulator
	font      font.Face
	smallFont font.Face
	buttons   []*Button
}

// NewUI creates a new UI instance
func NewUI(sim *simulator.Simulator) *UI {
	// Use the default font
	fontFace := text.FaceWithLineHeight(nil, 24)
	smallFontFace := text.FaceWithLineHeight(nil, 16)

	ui := &UI{
		simulator: sim,
		font:      fontFace,
		smallFont: smallFontFace,
	}

	// Create buttons
	ui.createButtons()

	return ui
}

// createButtons creates the UI buttons with a layout similar to the provided image
func (u *UI) createButtons() {
	buttonWidth := 100
	buttonHeight := 30
	buttonSpacing := 10
	startX := 20

	white := color.RGBA{255, 255, 255, 255}
	// Define button colors (matching the image)
	purple := color.RGBA{128, 0, 128, 255}
	gray := color.RGBA{128, 128, 128, 255}
	brown := color.RGBA{165, 42, 42, 255}
	blue := color.RGBA{65, 105, 225, 255}
	orange := color.RGBA{255, 140, 0, 255}
	green := color.RGBA{0, 128, 0, 255}
	red := color.RGBA{255, 0, 0, 255}

	u.buttons = []*Button{}

	// Row 1: View/Edit Mode buttons
	startY1 := 120
	u.buttons = append(u.buttons, &Button{
		X: startX, Y: startY1,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Reset View",
		BgColor: purple, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + buttonWidth + buttonSpacing, Y: startY1,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Grid",
		BgColor: gray, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*2, Y: startY1,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Snap",
		BgColor: gray, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*3, Y: startY1,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Edit Mode",
		BgColor: gray, TextColor: white,
	})

	// Row 2: Graph Operation buttons
	startY2 := startY1 + buttonHeight + buttonSpacing
	u.buttons = append(u.buttons, &Button{
		X: startX, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "New Graph",
		BgColor: purple, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + buttonWidth + buttonSpacing, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Load",
		BgColor: brown, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*2, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Save",
		BgColor: blue, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*3, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Add Edge",
		BgColor: blue, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*4, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Del Edge",
		BgColor: orange, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*5, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Add Node",
		BgColor: green, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*6, Y: startY2,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Del Node",
		BgColor: red, TextColor: white,
	})

	// Row 3: Algorithm/Simulation Control buttons
	startY3 := startY2 + buttonHeight + buttonSpacing
	u.buttons = append(u.buttons, &Button{
		X: startX, Y: startY3,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "BFS",
		BgColor: blue, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + buttonWidth + buttonSpacing, Y: startY3,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "DFS",
		BgColor: blue, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*2, Y: startY3,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "AVL Tree",
		BgColor: blue, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*3, Y: startY3,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Step",
		BgColor: green, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*4, Y: startY3,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Auto",
		BgColor: orange, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*5, Y: startY3,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Reset",
		BgColor: red, TextColor: white,
	})

	// Row 4: Advanced Graph Algorithms
	startY4 := startY3 + buttonHeight + buttonSpacing
	u.buttons = append(u.buttons, &Button{
		X: startX, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Dijkstra",
		BgColor: purple, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + buttonWidth + buttonSpacing, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "A*",
		BgColor: purple, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*2, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Topo Sort",
		BgColor: blue, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*3, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Kruskal",
		BgColor: green, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*4, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Prim",
		BgColor: green, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*5, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Tarjan",
		BgColor: orange, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*6, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Kosaraju",
		BgColor: orange, TextColor: white,
	})

	// Row 5: AVL Operation buttons
	startY5 := startY4 + buttonHeight + buttonSpacing
	u.buttons = append(u.buttons, &Button{
		X: startX, Y: startY5,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Insert",
		BgColor: green, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + buttonWidth + buttonSpacing, Y: startY5,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Delete",
		BgColor: red, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*2, Y: startY5,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Search",
		BgColor: orange, TextColor: white,
	})
}

// Draw draws the current state of the simulator
func (u *UI) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(color.RGBA{240, 240, 240, 255})

	// Draw title
	title := "Graph Algorithm Simulator"
	switch u.simulator.GetMode() {
	case algorithms.ModeAVL:
		title = "AVL Tree Simulator"
	case algorithms.ModeDijkstra:
		title = "Dijkstra's Shortest Path"
	case algorithms.ModeAStar:
		title = "A* Search Algorithm"
	case algorithms.ModeTopological:
		title = "Topological Sort"
	case algorithms.ModeKruskal:
		title = "Kruskal's MST"
	case algorithms.ModePrim:
		title = "Prim's MST"
	case algorithms.ModeTarjan:
		title = "Tarjan's SCC"
	case algorithms.ModeKosaraju:
		title = "Kosaraju's SCC"
	}
	text.Draw(screen, title, u.font, 20, 30, color.Black)

	// Draw buttons
	for _, button := range u.buttons {
		button.Draw(screen, nil)
	}

	// Draw current mode
	modeText := "Mode: "
	switch u.simulator.GetMode() {
	case algorithms.ModeBFS:
		modeText += "BFS"
	case algorithms.ModeDFS:
		modeText += "DFS"
	case algorithms.ModeAVL:
		modeText += "AVL Tree"
	case algorithms.ModeDijkstra:
		modeText += "Dijkstra"
	case algorithms.ModeAStar:
		modeText += "A*"
	case algorithms.ModeTopological:
		modeText += "Topological Sort"
	case algorithms.ModeKruskal:
		modeText += "Kruskal MST"
	case algorithms.ModePrim:
		modeText += "Prim MST"
	case algorithms.ModeTarjan:
		modeText += "Tarjan SCC"
	case algorithms.ModeKosaraju:
		modeText += "Kosaraju SCC"
	case algorithms.ModeIdle:
		modeText += "Idle"
	}
	text.Draw(screen, modeText, u.font, 20, 90, color.Black)

	// Draw AVL tree if in AVL mode
	if u.simulator.GetMode() == algorithms.ModeAVL {
		u.drawAVLTree(screen)
	} else {
		u.drawGraph(screen)
		u.drawAlgorithmResults(screen)
	}
}

// drawGraph draws the graph visualization
func (u *UI) drawGraph(screen *ebiten.Image) {
	// Draw edges
	for _, node := range u.simulator.Graph.Nodes {
		for j, neighbor := range node.Neighbors {
			// Draw edge
			ebitenutil.DrawLine(screen,
				float64(node.X), float64(node.Y),
				float64(u.simulator.Graph.Nodes[neighbor].X),
				float64(u.simulator.Graph.Nodes[neighbor].Y),
				color.Black)

			// Draw edge weight if it exists and we're in a weighted algorithm mode
			mode := u.simulator.GetMode()
			if (mode == algorithms.ModeDijkstra || mode == algorithms.ModeAStar ||
				mode == algorithms.ModeKruskal || mode == algorithms.ModePrim) &&
				j < len(node.Weights) {

				// Calculate midpoint
				midX := (node.X + u.simulator.Graph.Nodes[neighbor].X) / 2
				midY := (node.Y + u.simulator.Graph.Nodes[neighbor].Y) / 2

				// Draw weight
				weightText := fmt.Sprintf("%.1f", node.Weights[j])
				text.Draw(screen, weightText, u.smallFont,
					midX, midY, color.RGBA{128, 128, 128, 255})
			}
		}
	}

	// Draw nodes
	for i, node := range u.simulator.Graph.Nodes {
		// Determine node color based on state
		nodeColor := color.RGBA{100, 149, 237, 255} // Cornflower blue
		if i == u.simulator.Current {
			nodeColor = color.RGBA{255, 69, 0, 255} // Red-orange for current node
		} else if u.simulator.Visited[i] {
			nodeColor = color.RGBA{50, 205, 50, 255} // Lime green for visited nodes
		}

		// Draw node circle
		ebitenutil.DrawCircle(screen, float64(node.X), float64(node.Y), 20, nodeColor)
		ebitenutil.DrawCircle(screen, float64(node.X), float64(node.Y), 20, color.Black)

		// Draw node number
		nodeText := strconv.Itoa(i)
		bounds := text.BoundString(u.font, nodeText)
		text.Draw(screen, nodeText, u.font,
			node.X-bounds.Dx()/2,
			node.Y+bounds.Dy()/2,
			color.White)
	}
}

// drawAVLTree draws the AVL tree
func (u *UI) drawAVLTree(screen *ebiten.Image) {
	if u.simulator.GetAVLTree() == nil {
		return
	}

	// Draw tree nodes and edges
	u.drawAVLNode(screen, u.simulator.GetAVLTree().Root)
}

// drawAVLNode recursively draws an AVL tree node and its children
func (u *UI) drawAVLNode(screen *ebiten.Image, node *algorithms.AVLNode) {
	if node == nil {
		return
	}

	// Draw edges to children
	if node.Left != nil {
		u.drawAVLEdge(screen, node, node.Left)
	}
	if node.Right != nil {
		u.drawAVLEdge(screen, node, node.Right)
	}

	// Draw node
	nodeColor := color.RGBA{100, 149, 237, 255} // Cornflower blue
	if u.simulator.GetAVLAction() == "search" && u.simulator.GetAVLValue() == node.Value {
		nodeColor = color.RGBA{255, 69, 0, 255} // Red-orange for found node
	}

	// Draw node circle
	ebitenutil.DrawCircle(screen, float64(node.Position.X), float64(node.Position.Y), 20, nodeColor)
	ebitenutil.DrawCircle(screen, float64(node.Position.X), float64(node.Position.Y), 20, color.Black)

	// Draw node value
	valueText := strconv.Itoa(node.Value)
	bounds := text.BoundString(u.font, valueText)
	text.Draw(screen, valueText, u.font,
		node.Position.X-bounds.Dx()/2,
		node.Position.Y+bounds.Dy()/2,
		color.White)

	// Draw height
	heightText := strconv.Itoa(node.Height)
	bounds = text.BoundString(u.smallFont, heightText)
	text.Draw(screen, heightText, u.smallFont,
		node.Position.X-bounds.Dx()/2,
		node.Position.Y+30,
		color.Black)

	// Recursively draw children
	u.drawAVLNode(screen, node.Left)
	u.drawAVLNode(screen, node.Right)
}

// drawAVLEdge draws an edge between two AVL tree nodes
func (u *UI) drawAVLEdge(screen *ebiten.Image, from, to *algorithms.AVLNode) {
	// Calculate edge points
	x1, y1 := float64(from.Position.X), float64(from.Position.Y)
	x2, y2 := float64(to.Position.X), float64(to.Position.Y)

	// Draw line
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.Black)
}

// drawAlgorithmResults draws algorithm-specific visualization overlays
func (u *UI) drawAlgorithmResults(screen *ebiten.Image) {
	mode := u.simulator.GetMode()

	switch mode {
	case algorithms.ModeDijkstra:
		u.drawDijkstraResults(screen)
	case algorithms.ModeAStar:
		u.drawAStarPath(screen)
	case algorithms.ModeTopological:
		u.drawTopologicalOrder(screen)
	case algorithms.ModeKruskal, algorithms.ModePrim:
		u.drawMST(screen)
	case algorithms.ModeTarjan, algorithms.ModeKosaraju:
		u.drawSCCs(screen)
	}
}

// drawDijkstraResults draws shortest path distances
func (u *UI) drawDijkstraResults(screen *ebiten.Image) {
	distances := u.simulator.GetShortestPaths()
	if distances == nil {
		return
	}

	// Draw distance labels next to nodes
	for i, node := range u.simulator.Graph.Nodes {
		if dist, exists := distances[i]; exists && dist != math.Inf(1) {
			distText := fmt.Sprintf("%.1f", dist)
			text.Draw(screen, distText, u.smallFont,
				node.X+25, node.Y-10, color.RGBA{255, 0, 0, 255})
		}
	}
}

// drawAStarPath draws the found path
func (u *UI) drawAStarPath(screen *ebiten.Image) {
	path := u.simulator.GetPath()
	if path == nil || len(path) < 2 {
		return
	}

	// Draw path edges in red
	for i := 0; i < len(path)-1; i++ {
		from := u.simulator.Graph.Nodes[path[i]]
		to := u.simulator.Graph.Nodes[path[i+1]]
		ebitenutil.DrawLine(screen,
			float64(from.X), float64(from.Y),
			float64(to.X), float64(to.Y),
			color.RGBA{255, 0, 0, 255})
	}
}

// drawTopologicalOrder draws the topological ordering
func (u *UI) drawTopologicalOrder(screen *ebiten.Image) {
	topOrder := u.simulator.GetTopologicalOrder()
	if topOrder == nil {
		return
	}

	// Draw order numbers next to nodes
	for order, nodeId := range topOrder {
		if nodeId < len(u.simulator.Graph.Nodes) {
			node := u.simulator.Graph.Nodes[nodeId]
			orderText := fmt.Sprintf("%d", order+1)
			text.Draw(screen, orderText, u.smallFont,
				node.X+25, node.Y+10, color.RGBA{0, 0, 255, 255})
		}
	}
}

// drawMST draws the minimum spanning tree edges
func (u *UI) drawMST(screen *ebiten.Image) {
	mst := u.simulator.GetMST()
	if mst == nil {
		return
	}

	// Draw MST edges in green
	for _, edge := range mst {
		if edge.From < len(u.simulator.Graph.Nodes) && edge.To < len(u.simulator.Graph.Nodes) {
			from := u.simulator.Graph.Nodes[edge.From]
			to := u.simulator.Graph.Nodes[edge.To]
			ebitenutil.DrawLine(screen,
				float64(from.X), float64(from.Y),
				float64(to.X), float64(to.Y),
				color.RGBA{0, 255, 0, 255})

			// Draw weight label
			midX := (from.X + to.X) / 2
			midY := (from.Y + to.Y) / 2
			weightText := fmt.Sprintf("%.1f", edge.Weight)
			text.Draw(screen, weightText, u.smallFont,
				midX, midY, color.RGBA{0, 128, 0, 255})
		}
	}
}

// drawSCCs draws strongly connected components with different colors
func (u *UI) drawSCCs(screen *ebiten.Image) {
	sccs := u.simulator.GetSCCs()
	if sccs == nil {
		return
	}

	// Define colors for different SCCs
	sccColors := []color.RGBA{
		{255, 100, 100, 255}, // Red
		{100, 255, 100, 255}, // Green
		{100, 100, 255, 255}, // Blue
		{255, 255, 100, 255}, // Yellow
		{255, 100, 255, 255}, // Magenta
		{100, 255, 255, 255}, // Cyan
	}

	// Draw SCC labels
	for sccIndex, scc := range sccs {
		colorIndex := sccIndex % len(sccColors)
		sccColor := sccColors[colorIndex]

		for _, nodeId := range scc {
			if nodeId < len(u.simulator.Graph.Nodes) {
				node := u.simulator.Graph.Nodes[nodeId]
				sccText := fmt.Sprintf("SCC%d", sccIndex+1)
				text.Draw(screen, sccText, u.smallFont,
					node.X-15, node.Y+30, sccColor)
			}
		}
	}
}

// Update updates the UI state
func (u *UI) Update() error {
	// Get mouse position
	x, y := ebiten.CursorPosition()

	// Update button hover states and handle clicks
	for i, button := range u.buttons {
		// Check if mouse is over button
		button.Hover = button.IsInside(x, y)

		// Handle button clicks
		if button.Hover && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			switch i {
			case 0: // Reset View
				// Implementation needed
			case 1: // Grid
				// Implementation needed
			case 2: // Snap
				// Implementation needed
			case 3: // Edit Mode
				// Implementation needed
			case 4: // New Graph
				// Implementation needed
			case 5: // Load
				// Implementation needed
			case 6: // Save
				// Implementation needed
			case 7: // Add Edge
				// Implementation needed
			case 8: // Del Edge
				// Implementation needed
			case 9: // Add Node
				// Implementation needed
			case 10: // Del Node
				// Implementation needed
			case 11: // BFS
				u.simulator.StartBFS(0)
			case 12: // DFS
				u.simulator.StartDFS(0)
			case 13: // AVL Tree
				u.simulator.StartAVL()
			case 14: // Step
				u.simulator.Update()
			case 15: // Auto
				// Implementation needed
			case 16: // Reset
				u.simulator.Reset()
			case 17: // Dijkstra
				u.simulator.StartDijkstra(0)
			case 18: // A*
				u.simulator.StartAStar(0, len(u.simulator.Graph.Nodes)-1) // From first to last node
			case 19: // Topological Sort
				u.simulator.StartTopological()
			case 20: // Kruskal
				u.simulator.StartKruskal()
			case 21: // Prim
				u.simulator.StartPrim()
			case 22: // Tarjan
				u.simulator.StartTarjan()
			case 23: // Kosaraju
				u.simulator.StartKosaraju()
			case 24: // Insert (AVL)
				if u.simulator.GetMode() == algorithms.ModeAVL {
					u.simulator.SetAVLAction("insert")
					u.simulator.InsertAVL(u.simulator.GetAVLValue())
					u.simulator.IncrementAVLValue()
				}
			case 25: // Delete (AVL)
				if u.simulator.GetMode() == algorithms.ModeAVL {
					u.simulator.SetAVLAction("delete")
					if u.simulator.GetAVLValue() > 0 {
						u.simulator.DecrementAVLValue()
						u.simulator.DeleteAVL(u.simulator.GetAVLValue())
					}
				}
			case 26: // Search (AVL)
				if u.simulator.GetMode() == algorithms.ModeAVL {
					u.simulator.SetAVLAction("search")
					u.simulator.SearchAVL(u.simulator.GetAVLValue())
				}
			}
		}
	}

	// Handle keyboard shortcuts for now (will be removed later)

	return nil
}
