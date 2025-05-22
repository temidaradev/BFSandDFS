package ui

import (
	"image/color"
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

	// Row 4: AVL Operation buttons
	startY4 := startY3 + buttonHeight + buttonSpacing
	u.buttons = append(u.buttons, &Button{
		X: startX, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Insert",
		BgColor: green, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + buttonWidth + buttonSpacing, Y: startY4,
		Width: buttonWidth, Height: buttonHeight,
		Text:    "Delete",
		BgColor: red, TextColor: white,
	})
	u.buttons = append(u.buttons, &Button{
		X: startX + (buttonWidth+buttonSpacing)*2, Y: startY4,
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
	title := "Graph Traversal Simulator"
	if u.simulator.GetMode() == algorithms.ModeAVL {
		title = "AVL Tree Simulator"
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
	case algorithms.ModeIdle:
		modeText += "Idle"
	}
	text.Draw(screen, modeText, u.font, 20, 90, color.Black)

	// Draw AVL tree if in AVL mode
	if u.simulator.GetMode() == algorithms.ModeAVL {
		u.drawAVLTree(screen)
	} else {
		u.drawGraph(screen)
	}
}

// drawGraph draws the graph visualization
func (u *UI) drawGraph(screen *ebiten.Image) {
	// Draw edges
	for _, node := range u.simulator.Graph.Nodes {
		for _, neighbor := range node.Neighbors {
			// Draw edge
			ebitenutil.DrawLine(screen,
				float64(node.X), float64(node.Y),
				float64(u.simulator.Graph.Nodes[neighbor].X),
				float64(u.simulator.Graph.Nodes[neighbor].Y),
				color.Black)
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
			case 17: // Insert (AVL)
				if u.simulator.GetMode() == algorithms.ModeAVL {
					u.simulator.SetAVLAction("insert")
					u.simulator.InsertAVL(u.simulator.GetAVLValue())
					u.simulator.IncrementAVLValue()
				}
			case 18: // Delete (AVL)
				if u.simulator.GetMode() == algorithms.ModeAVL {
					u.simulator.SetAVLAction("delete")
					if u.simulator.GetAVLValue() > 0 {
						u.simulator.DecrementAVLValue()
						u.simulator.DeleteAVL(u.simulator.GetAVLValue())
					}
				}
			case 19: // Search (AVL)
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
