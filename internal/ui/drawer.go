package ui

import (
	"fmt"
	"image/color"
	"strings"

	"bfsdfs/internal/algorithms"
	"bfsdfs/pkg/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/temidaradev/esset/v2"
	"github.com/temidaradev/esset/v2/example/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

func init() {
	assets.FontFaceS, _ = esset.GetFont(assets.MyFont, 12)
}

func MeasureText(face text.Face, s string) (width, height float64) {
	width, height = text.Measure(s, face, 0)
	return
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Get window dimensions
	screenWidth, screenHeight := ebiten.WindowSize()

	// Only redraw if necessary
	if g.canvasNeedsRedraw {
		// Create a separate canvas for graph visualization - use full screen instead of reserving space
		if g.graphCanvas == nil || g.graphCanvas.Bounds().Dx() != screenWidth || g.graphCanvas.Bounds().Dy() != screenHeight {
			g.graphCanvas = ebiten.NewImage(screenWidth, screenHeight)
		}
		g.graphCanvas.Fill(color.RGBA{240, 240, 240, 255})

		// Check if we're in AVL mode and draw accordingly
		if g.Sim.Mode == algorithms.ModeAVL {
			// Draw AVL tree
			g.drawAVLTree(g.graphCanvas)
		} else {
			// Draw normal graph
			g.drawGraph(g.graphCanvas, screenWidth, screenHeight)
		}

		g.canvasNeedsRedraw = false
	}

	// Draw the cached graph canvas
	screen.DrawImage(g.graphCanvas, nil)

	// Draw selection box if selecting
	if g.Selecting {
		// Determine the boundaries of the selection box in screen coordinates
		left := min(g.SelectionStartX, g.MouseX)
		right := max(g.SelectionStartX, g.MouseX)
		top := min(g.SelectionStartY, g.MouseY)
		bottom := max(g.SelectionStartY, g.MouseY)

		// Draw a transparent rectangle for the selection area
		selectionColor := color.RGBA{100, 150, 200, 50} // Light blue with transparency
		draw.DrawRect(screen, float64(left), float64(top), float64(right-left), float64(bottom-top), selectionColor)

		// Draw a border around the selection area
		borderColor := color.RGBA{100, 150, 200, 255}                                                       // Opaque light blue
		draw.DrawLine(screen, float64(left), float64(top), float64(right), float64(top), borderColor)       // Top border
		draw.DrawLine(screen, float64(left), float64(bottom), float64(right), float64(bottom), borderColor) // Bottom border
		draw.DrawLine(screen, float64(left), float64(top), float64(left), float64(bottom), borderColor)     // Left border
		draw.DrawLine(screen, float64(right), float64(top), float64(right), float64(bottom), borderColor)   // Right border
	}

	// Draw the button window
	titleBarHeight := 25
	var windowWidth int
	if g.ButtonsCollapsed {
		windowWidth = 120
	} else {
		windowWidth = g.ButtonWindowWidth
	}

	var windowHeight int
	if g.ButtonsCollapsed {
		windowHeight = titleBarHeight
	} else {
		windowHeight = g.ButtonWindowHeight
	}
	windowBg := ebiten.NewImage(windowWidth, windowHeight)

	// Window background
	windowBg.Fill(color.RGBA{220, 220, 220, 240}) // Light gray with slight transparency

	// Title bar background
	titleBar := ebiten.NewImage(windowWidth, titleBarHeight)
	titleBar.Fill(color.RGBA{70, 130, 180, 255}) // Steel blue title bar

	// Apply window
	windowOpts := &ebiten.DrawImageOptions{}
	windowOpts.GeoM.Translate(float64(g.ButtonWindowX), float64(g.ButtonWindowY))
	screen.DrawImage(windowBg, windowOpts)

	// Apply title bar
	titleBarOpts := &ebiten.DrawImageOptions{}
	titleBarOpts.GeoM.Translate(float64(g.ButtonWindowX), float64(g.ButtonWindowY))
	screen.DrawImage(titleBar, titleBarOpts)

	// Draw window title
	titleText := "Graph Controls"
	//text.Draw(screen, titleText, basicfont.Face7x13, g.ButtonWindowX+10, g.ButtonWindowY+16, color.White)
	esset.DrawText(screen, titleText, float64(g.ButtonWindowX+10), float64(g.ButtonWindowY+5), assets.FontFaceS, color.White)

	// Draw collapse button in title bar
	collapseButtonWidth := 25
	collapseButtonHeight := 20
	collapseButtonX := g.ButtonWindowX + windowWidth - collapseButtonWidth - 5
	collapseButtonY := g.ButtonWindowY + 3

	collapseButtonBg := ebiten.NewImage(collapseButtonWidth, collapseButtonHeight)
	collapseButtonBg.Fill(color.RGBA{90, 150, 200, 255})

	collapseButtonOpts := &ebiten.DrawImageOptions{}
	collapseButtonOpts.GeoM.Translate(float64(collapseButtonX), float64(collapseButtonY))
	screen.DrawImage(collapseButtonBg, collapseButtonOpts)

	// Draw button text (centered)
	var collapseText string
	if g.ButtonsCollapsed {
		collapseText = ">>"
	} else {
		collapseText = "<<"
	}
	textX := collapseButtonX + (collapseButtonWidth-len(collapseText)*7)/2 // Approximate width
	textY := collapseButtonY + 2
	//text.Draw(screen, collapseText, basicfont.Face7x13, textX, textY, color.White)
	esset.DrawText(screen, collapseText, float64(textX), float64(textY), assets.FontFaceS, color.White)

	// Store collapse button bounds for click detection
	g.CollapseButton.X = collapseButtonX
	g.CollapseButton.Y = collapseButtonY
	g.CollapseButton.Width = collapseButtonWidth
	g.CollapseButton.Height = collapseButtonHeight
	g.CollapseButton.Text = collapseText

	// Draw buttons only if window is not collapsed
	if !g.ButtonsCollapsed {
		// Initialize the grid layout
		titleBarHeight := 25
		buttonOffsetY := titleBarHeight + 10 // Start below title bar
		buttonOffsetX := 20                  // Left margin inside window

		// Calculate available space
		availableWidth := g.ButtonWindowWidth - (buttonOffsetX * 2) // Subtract left and right margins

		// Get the widest button to calculate proper spacing
		maxButtonWidth := 0
		for _, btn := range g.Buttons {
			if btn.Width > maxButtonWidth {
				maxButtonWidth = btn.Width
			}
		}

		// Calculate how many buttons can fit per row
		buttonMargin := 10
		buttonsPerRow := max(1, (availableWidth)/(maxButtonWidth+buttonMargin))

		// Standard button height (used by all buttons)
		buttonHeight := 30

		// Draw buttons in a grid layout
		for i, btn := range g.Buttons {
			// Store original values to restore them later
			origX := btn.X
			origY := btn.Y
			origAnchorBottom := btn.AnchorBottom

			// Calculate position in the grid
			rowIndex := i / buttonsPerRow
			colIndex := i % buttonsPerRow

			// Adjust X position to ensure buttons are evenly spaced
			cellWidth := availableWidth / buttonsPerRow
			btnXOffset := (cellWidth - btn.Width) / 2

			// Set position relative to window
			btn.X = g.ButtonWindowX + buttonOffsetX + (colIndex * cellWidth) + btnXOffset
			btn.Y = g.ButtonWindowY + buttonOffsetY + rowIndex*(btn.Height+buttonMargin)
			btn.AnchorBottom = false // Disable bottom anchoring

			// Draw the button
			btn.Draw(screen, nil)

			// Restore original values so we don't affect the button's original position
			btn.X = origX
			btn.Y = origY
			btn.AnchorBottom = origAnchorBottom
		}

		// Update window height to fit all buttons
		buttonCount := len(g.Buttons)
		rowCount := (buttonCount + buttonsPerRow - 1) / buttonsPerRow                  // Ceiling division
		totalButtonHeight := rowCount*(buttonHeight+buttonMargin) + buttonOffsetY + 10 // +10 for bottom margin

		// Update window height if needed (only when not collapsed)
		if !g.ButtonsCollapsed {
			g.ButtonWindowHeight = max(titleBarHeight+50, totalButtonHeight)
		}
	}

	// Draw algorithm info if active (Visit order, Queue/Stack)
	if g.Sim.Mode != algorithms.ModeIdle && g.Sim.Mode != algorithms.ModeAVL {
		// Draw visit order
		orderStr := "Visit order: "
		for i, nodeIdx := range g.Sim.Order {
			if i > 0 {
				orderStr += " > "
			}
			orderStr += string(rune('A' + nodeIdx))
		}
		// Position visit order at the top, slightly below the screen edge
		//text.Draw(screen, orderStr, basicfont.Face7x13, 20, 20, color.Black)
		esset.DrawText(screen, orderStr, 20, 10, assets.FontFaceS, color.Black)

		// Draw queue or stack status
		var dataStructStr string
		if g.Sim.Mode == algorithms.ModeBFS {
			dataStructStr = "Queue: "
			for i, nodeIdx := range g.Sim.Queue {
				if i > 0 {
					dataStructStr += ", "
				}
				dataStructStr += string(rune('A' + nodeIdx))
			}
		} else if g.Sim.Mode == algorithms.ModeDFS {
			dataStructStr = "Stack: "
			for i, nodeIdx := range g.Sim.Stack {
				if i > 0 {
					dataStructStr += ", "
				}
				dataStructStr += string(rune('A' + nodeIdx))
			}
		}
		// Position queue/stack status below visit order
		//text.Draw(screen, dataStructStr, basicfont.Face7x13, 20, 40, color.Black)
		esset.DrawText(screen, dataStructStr, 20, 30, assets.FontFaceS, color.Black)
	} else if g.Sim.Mode == algorithms.ModeAVL {
		// Draw AVL tree info
		avlInfoStr := "AVL Tree Mode"
		//text.Draw(screen, avlInfoStr, basicfont.Face7x13, 20, 20, color.Black)
		esset.DrawText(screen, avlInfoStr, 20, 10, assets.FontFaceS, color.Black)

		// Show current action if any
		if g.Sim.GetAVLAction() != "" {
			actionStr := fmt.Sprintf("Last Action: %s", g.Sim.GetAVLAction())
			//text.Draw(screen, actionStr, basicfont.Face7x13, 20, 40, color.Black)
			esset.DrawText(screen, actionStr, 20, 30, assets.FontFaceS, color.Black)
		}

		// Show tree statistics
		if g.Sim.GetAVLTree() != nil && g.Sim.GetAVLTree().Root != nil {
			heightStr := fmt.Sprintf("Tree Height: %d", g.Sim.GetAVLTree().Root.Height)
			//text.Draw(screen, heightStr, basicfont.Face7x13, 20, 60, color.Black)
			esset.DrawText(screen, heightStr, 20, 50, assets.FontFaceS, color.Black)
		}
	}

	// Draw the message display
	if g.MessageTimer > 0 {
		// Background for message
		messageBgWidth := 300 // Fixed width for the message box
		messageBgHeight := 30
		messageBgX := (screenWidth - messageBgWidth) / 2 // Center horizontally
		// Position message near the bottom, above the speed slider and zoom indicator
		messageBgY := screenHeight - 80 // Adjust position
		messageBg := ebiten.NewImage(messageBgWidth, messageBgHeight)
		messageBg.Fill(color.RGBA{50, 50, 50, 200}) // Dark gray with transparency
		messageOpts := &ebiten.DrawImageOptions{}
		messageOpts.GeoM.Translate(float64(messageBgX), float64(messageBgY))
		screen.DrawImage(messageBg, messageOpts)

		messageText := g.Message

		messageWidth, messageHeight := MeasureText(assets.FontFaceS, messageText)

		messageTextX := float64(messageBgX) + (float64(messageBgWidth)-messageWidth)/2
		messageTextY := float64(messageBgY) + (float64(messageBgHeight)-messageHeight)/2

		esset.DrawText(screen, messageText, messageTextX, messageTextY, assets.FontFaceS, color.White)

	}

	// Draw speed slider
	// Background
	sliderBgWidth := 200
	sliderBgHeight := 20
	// Position speed slider near the bottom right
	sliderBgX := screenWidth - sliderBgWidth - 20
	sliderBgY := screenHeight - 30 // Adjust position
	sliderBg := ebiten.NewImage(sliderBgWidth, sliderBgHeight)
	sliderBg.Fill(color.RGBA{80, 80, 80, 255})
	sliderOpts := &ebiten.DrawImageOptions{}
	sliderOpts.GeoM.Translate(float64(sliderBgX), float64(sliderBgY))
	screen.DrawImage(sliderBg, sliderOpts)

	// Handle
	handleWidth := 10
	handleHeight := 20
	// Calculate handle position based on StepDelay (50 to 10)
	// StepDelay 50 = left side (slow), StepDelay 10 = right side (fast)
	handleX := sliderBgX + int(float64(50-g.StepDelay)/40.0*float64(sliderBgWidth-handleWidth))
	handleY := sliderBgY
	handle := ebiten.NewImage(handleWidth, handleHeight)
	handle.Fill(color.RGBA{200, 200, 200, 255})
	handleOpts := &ebiten.DrawImageOptions{}
	handleOpts.GeoM.Translate(float64(handleX), float64(handleY))
	screen.DrawImage(handle, handleOpts)

	// Speed label
	speedLabel := fmt.Sprintf("Speed: %d", 50-g.StepDelay+10)
	labelWidth, labelHeight := MeasureText(assets.FontFaceS, speedLabel)

	x := float64(sliderBgX) - labelWidth - 10
	y := float64(sliderBgY) + labelHeight/2 - 6 // Rough vertical center

	esset.DrawText(screen, speedLabel, x, y, assets.FontFaceS, color.Black)

	// Draw Help Overlay
	if g.ShowHelp {
		drawHelpOverlay(screen, screenWidth, screenHeight)
	}

	// Draw Context Menu
	if g.ContextMenu.Visible {
		// ContextMenu.Draw likely only needs the screen
		g.ContextMenu.Draw(screen)
	}

	// Draw File Dialogs
	if g.ShowSaveDialog {
		// FileDialog.Draw likely only needs the screen
		g.SaveDialog.Draw(screen)
	}
	if g.ShowLoadDialog {
		// FileDialog.Draw likely only needs the screen
		g.LoadDialog.Draw(screen)
	}

	// Draw AVL Input Modal
	if g.ShowAVLInput {
		// Dim the background
		dimming := ebiten.NewImage(screenWidth, screenHeight)
		dimming.Fill(color.RGBA{0, 0, 0, 100}) // Semi-transparent black
		screen.DrawImage(dimming, nil)

		// Modal background
		modalWidth := 300
		modalHeight := 150
		modalX := (screenWidth - modalWidth) / 2
		modalY := (screenHeight - modalHeight) / 2
		modalBg := ebiten.NewImage(modalWidth, modalHeight)
		modalBg.Fill(color.RGBA{200, 200, 200, 255}) // Light gray
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(modalX), float64(modalY))
		screen.DrawImage(modalBg, opts)

		// Modal title
		title := fmt.Sprintf("%s Value", strings.Title(g.AVLAction))
		//text.Draw(screen, title, basicfont.Face7x13, modalX+10, modalY+20, color.Black)
		esset.DrawText(screen, title, float64(modalX+10), float64(modalY+10), assets.FontFaceS, color.Black)

		// Input field background
		inputWidth := 280
		inputHeight := 30
		inputX := modalX + 10
		inputY := modalY + 40
		inputBg := ebiten.NewImage(inputWidth, inputHeight)
		inputBg.Fill(color.White)
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(inputX), float64(inputY))
		screen.DrawImage(inputBg, opts)

		// Input field text
		//text.Draw(screen, g.AVLInputText, basicfont.Face7x13, inputX+5, inputY+inputHeight/2+basicfont.Face7x13.Ascent/2, color.Black)
		esset.DrawText(screen, g.AVLInputText, float64(inputX+5), float64(inputY+inputHeight/2+basicfont.Face7x13.Ascent/2-10), assets.FontFaceS, color.Black)

		// Action buttons
		buttonWidth := 80
		buttonHeight := 30
		buttonSpacing := 10
		buttonY := modalY + modalHeight - buttonHeight - 10

		// OK button
		okButtonX := modalX + modalWidth - buttonWidth*2 - buttonSpacing*2
		drawButton(screen, okButtonX, buttonY, buttonWidth, buttonHeight, "OK", color.RGBA{100, 150, 100, 255}, color.RGBA{255, 255, 255, 255}, basicfont.Face7x13)

		// Cancel button
		cancelButtonX := modalX + modalWidth - buttonWidth - buttonSpacing
		drawButton(screen, cancelButtonX, buttonY, buttonWidth, buttonHeight, "Cancel", color.RGBA{150, 100, 100, 255}, color.RGBA{255, 255, 255, 255}, basicfont.Face7x13)
	}
}

// drawButton is a helper function to draw a button
func drawButton(screen *ebiten.Image, x, y, width, height int, textLabel string, bgColor, textColor color.RGBA, face font.Face) {
	// Draw button background
	buttonImage := ebiten.NewImage(width, height)
	buttonImage.Fill(bgColor)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(buttonImage, opts)

	// Draw button text
	labelWidth, labelHeight := text.Measure(textLabel, assets.FontFaceS, 0)

	textX := float64(x) + (float64(width)-labelWidth)/2
	textY := float64(y) + (float64(height)-labelHeight)/2

	esset.DrawText(screen, textLabel, textX, textY, assets.FontFaceS, textColor)
}

// drawHelpOverlay draws the help information overlay
func drawHelpOverlay(screen *ebiten.Image, screenWidth, screenHeight int) {
	// Dim the background
	dimming := ebiten.NewImage(screenWidth, screenHeight)
	dimming.Fill(color.RGBA{0, 0, 0, 150}) // Semi-transparent black
	screen.DrawImage(dimming, nil)

	// Help text background
	helpBgWidth := 400
	helpBgHeight := 300
	helpBgX := (screenWidth - helpBgWidth) / 2
	helpBgY := (screenHeight - helpBgHeight) / 2
	helpBg := ebiten.NewImage(helpBgWidth, helpBgHeight)
	helpBg.Fill(color.RGBA{220, 220, 220, 255}) // Light gray
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(helpBgX), float64(helpBgY))
	screen.DrawImage(helpBg, opts)

	// Help text content
	helpText := `Keyboard Shortcuts:

BFS/DFS Simulation:
  SPACE: Step forward
  A: Toggle Auto-step
  R: Reset simulation

Graph Editing:
  E: Toggle Edit Mode
  N: Add Node (click on canvas)
  D: Delete Node (click on node)
  A: Add Edge (click two nodes)
  X: Delete Edge (click two nodes)

View Controls:
  Middle Click / Shift+Right Click: Pan
  H: Toggle Help

Context Menu:
  Right Click on node or empty area
`
	esset.DrawText(screen, helpText, float64(helpBgX+20), float64(helpBgY+20), assets.FontFaceS, color.Black)

	// Draw close instruction (bottom-right, 20px padding)
	closeText := "Press H to close"
	closeWidth, _ := text.Measure(closeText, assets.FontFaceS, 0)

	closeTextX := float64(helpBgX+helpBgWidth) - closeWidth - 20
	closeTextY := float64(helpBgY+helpBgHeight) - 20 // already bottom padding

	esset.DrawText(screen, closeText, closeTextX, closeTextY, assets.FontFaceS, color.Black)

}

// drawGraph draws the normal graph visualization
func (g *Game) drawGraph(canvas *ebiten.Image, screenWidth, screenHeight int) {
	// Draw grid if enabled - draw it on the graph canvas so it moves with the graph
	if g.ShowGrid {
		// Create a temporary canvas for the grid that's larger than the screen
		// to ensure we have enough grid lines when moving
		gridSize := 1000 // Make it large enough to cover movement range
		if g.gridCanvas == nil || g.gridCanvas.Bounds().Dx() != gridSize || g.gridCanvas.Bounds().Dy() != gridSize {
			g.gridCanvas = ebiten.NewImage(gridSize, gridSize)
		}
		g.gridCanvas.Fill(color.RGBA{240, 240, 240, 255})

		// Draw grid on the temporary canvas using optimized drawing
		draw.DrawOptimizedGrid(g.gridCanvas, gridSize, gridSize, g.GridConfig)

		// Draw grid border
		borderColor := color.RGBA{100, 100, 100, 255}
		// Draw top and bottom borders
		for i := 0; i < gridSize; i++ {
			g.gridCanvas.Set(i, 0, borderColor)
			g.gridCanvas.Set(i, gridSize-1, borderColor)
		}
		// Draw left and right borders
		for i := 0; i < gridSize; i++ {
			g.gridCanvas.Set(0, i, borderColor)
			g.gridCanvas.Set(gridSize-1, i, borderColor)
		}

		// Draw the grid canvas onto the graph canvas with offset
		gridOpts := &ebiten.DrawImageOptions{}
		gridOpts.GeoM.Translate(g.CanvasOffsetX, g.CanvasOffsetY)
		canvas.DrawImage(g.gridCanvas, gridOpts)
	}

	// Draw edges
	for _, edge := range g.Sim.Graph.Edges {
		// Get node positions
		node1 := g.Sim.Graph.Nodes[edge[0]]
		node2 := g.Sim.Graph.Nodes[edge[1]]

		// Convert node positions to screen coordinates
		x1 := float64(node1.X) + g.CanvasOffsetX
		y1 := float64(node1.Y) + g.CanvasOffsetY
		x2 := float64(node2.X) + g.CanvasOffsetX
		y2 := float64(node2.Y) + g.CanvasOffsetY

		// Check if edge is visible on screen
		if x1 < float64(screenWidth) && x2 < float64(screenWidth) &&
			x1 > 0 && x2 > 0 &&
			y1 < float64(screenHeight) && y2 < float64(screenHeight) &&
			y1 > 0 && y2 > 0 {

			// Draw edge
			edgeColor := color.RGBA{100, 100, 100, 255}
			draw.DrawCachedLine(canvas, x1, y1, x2, y2, edgeColor)
		}
	}

	// Draw nodes
	for i, node := range g.Sim.Graph.Nodes {
		// Convert node position to screen coordinates
		x := float64(node.X) + g.CanvasOffsetX
		y := float64(node.Y) + g.CanvasOffsetY

		// Check if node is visible on screen
		if x < float64(screenWidth) && x > 0 && y < float64(screenHeight) && y > 0 {
			// Determine node color based on state
			var nodeColor color.RGBA
			if i == g.Sim.Current {
				nodeColor = color.RGBA{255, 69, 0, 255} // Red-orange for current node
			} else if g.Sim.Visited[i] {
				nodeColor = color.RGBA{50, 205, 50, 255} // Lime green for visited nodes
			} else {
				nodeColor = color.RGBA{70, 130, 180, 255} // Cornflower blue for unvisited nodes
			}

			// Draw node (fixed radius of 20)
			draw.DrawCachedCircle(canvas, int(x), int(y), 20, nodeColor)

			// Draw node label
			label := string(rune('A' + i))
			//text.Draw(canvas, label, basicfont.Face7x13, int(x)-3, int(y)+4, color.White)
			esset.DrawText(canvas, label, float64(int(x)-4), float64(int(y)-7), assets.FontFaceS, color.White)
		}
	}
}

// drawAVLTree draws the AVL tree visualization
func (g *Game) drawAVLTree(canvas *ebiten.Image) {
	if g.Sim.GetAVLTree() == nil || g.Sim.GetAVLTree().Root == nil {
		return
	}

	// Update node positions for visualization
	screenWidth, _ := ebiten.WindowSize()
	centerX := screenWidth / 2
	startY := 100
	levelHeight := 80

	// Update positions
	g.Sim.GetAVLTree().UpdatePositions(centerX, startY, levelHeight)

	// Draw tree nodes and edges
	g.drawAVLNode(canvas, g.Sim.GetAVLTree().Root)
}

// drawAVLNode recursively draws an AVL tree node and its children
func (g *Game) drawAVLNode(canvas *ebiten.Image, node *algorithms.AVLNode) {
	if node == nil {
		return
	}

	// Draw edges to children first (so they appear behind nodes)
	if node.Left != nil {
		g.drawAVLEdge(canvas, node, node.Left)
		g.drawAVLNode(canvas, node.Left)
	}
	if node.Right != nil {
		g.drawAVLEdge(canvas, node, node.Right)
		g.drawAVLNode(canvas, node.Right)
	}

	// Draw node
	nodeColor := color.RGBA{100, 149, 237, 255} // Cornflower blue
	if g.Sim.GetAVLAction() == "search" && g.Sim.GetAVLValue() == node.Value {
		nodeColor = color.RGBA{255, 69, 0, 255} // Red-orange for found node
	}

	// Apply canvas offset (no zoom scaling)
	x := float64(node.Position.X) + g.CanvasOffsetX
	y := float64(node.Position.Y) + g.CanvasOffsetY

	// Draw node circle with border (fixed radius of 25)
	draw.DrawCachedCircle(canvas, int(x), int(y), 25, nodeColor)
	draw.DrawCachedCircle(canvas, int(x), int(y), 25, color.RGBA{0, 0, 0, 255}) // Black border

	// Draw node value
	valueText := fmt.Sprintf("%d", node.Value)
	valueWidth, valueHeight := text.Measure(valueText, assets.FontFaceS, 0)

	// Center horizontally at x, vertically align baseline approximately at y
	valueX := float64(x) - valueWidth/2
	valueY := float64(y) + valueHeight*0.8 // baseline approximation

	esset.DrawText(canvas, valueText, valueX, valueY, assets.FontFaceS, color.White)

	// Draw height below the node
	heightText := fmt.Sprintf("h:%d", node.Height)
	heightWidth, heightHeight := text.Measure(heightText, assets.FontFaceS, 0)

	heightX := float64(x) - heightWidth/2
	heightY := float64(y) + 35 + heightHeight*0.8 // 35 px below + baseline approx

	esset.DrawText(canvas, heightText, heightX, heightY, assets.FontFaceS, color.Black)
}

// drawAVLEdge draws an edge between two AVL tree nodes
func (g *Game) drawAVLEdge(canvas *ebiten.Image, from, to *algorithms.AVLNode) {
	// Apply canvas offset (no zoom scaling)
	x1 := float64(from.Position.X) + g.CanvasOffsetX
	y1 := float64(from.Position.Y) + g.CanvasOffsetY
	x2 := float64(to.Position.X) + g.CanvasOffsetX
	y2 := float64(to.Position.Y) + g.CanvasOffsetY

	// Draw line
	draw.DrawCachedLine(canvas, x1, y1, x2, y2, color.RGBA{0, 0, 0, 255})
}
