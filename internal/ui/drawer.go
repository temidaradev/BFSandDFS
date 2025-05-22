package ui

import (
	"fmt"
	"image/color"
	"strings"

	"bfsdfs/internal/algorithms"
	"bfsdfs/pkg/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

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
			g.graphCanvas.DrawImage(g.gridCanvas, gridOpts)
		}

		// Draw edges
		for _, edge := range g.Sim.Graph.Edges {
			// Get node positions
			node1 := g.Sim.Graph.Nodes[edge[0]]
			node2 := g.Sim.Graph.Nodes[edge[1]]

			// Convert node positions to screen coordinates
			x1 := float64(node1.X)*g.ZoomLevel + g.CanvasOffsetX
			y1 := float64(node1.Y)*g.ZoomLevel + g.CanvasOffsetY
			x2 := float64(node2.X)*g.ZoomLevel + g.CanvasOffsetX
			y2 := float64(node2.Y)*g.ZoomLevel + g.CanvasOffsetY

			// Check if edge is visible on screen
			if x1 < float64(screenWidth) && x2 < float64(screenWidth) &&
				x1 > 0 && x2 > 0 &&
				y1 < float64(screenHeight) && y2 < float64(screenHeight) &&
				y1 > 0 && y2 > 0 {

				// Draw edge
				edgeColor := color.RGBA{100, 100, 100, 255}
				draw.DrawCachedLine(g.graphCanvas, x1, y1, x2, y2, edgeColor)
			}
		}

		// Draw nodes
		for i, node := range g.Sim.Graph.Nodes {
			// Convert node position to screen coordinates
			x := float64(node.X)*g.ZoomLevel + g.CanvasOffsetX
			y := float64(node.Y)*g.ZoomLevel + g.CanvasOffsetY

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

				// Draw node
				draw.DrawCachedCircle(g.graphCanvas, int(x), int(y), int(20*g.ZoomLevel), nodeColor)

				// Draw node label
				label := string(rune('A' + i))
				text.Draw(g.graphCanvas, label, basicfont.Face7x13, int(x)-3, int(y)+4, color.White)
			}
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

	// Draw HUD overlay (no black bars, just elements)
	// topHudHeight := 40 // Removed top HUD
	// bottomHudHeight := 60 // Removed bottom HUD

	// Draw buttons directly
	for _, btn := range g.Buttons {
		btn.Draw(screen, g)
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
		text.Draw(screen, orderStr, basicfont.Face7x13, 20, 20, color.Black)

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
		text.Draw(screen, dataStructStr, basicfont.Face7x13, 20, 40, color.Black)
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

		// Message text
		messageText := g.Message
		// Center the text within the background
		messageBounds := text.BoundString(basicfont.Face7x13, messageText)
		messageTextX := messageBgX + (messageBgWidth-messageBounds.Dx())/2
		messageTextY := messageBgY + (messageBgHeight-messageBounds.Dy())/2 + basicfont.Face7x13.Ascent
		text.Draw(screen, messageText, basicfont.Face7x13, messageTextX, messageTextY, color.White)
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
	handleX := sliderBgX + int(float64(50-g.StepDelay)/40.0*float64(sliderBgWidth-handleWidth))
	handleY := sliderBgY
	handle := ebiten.NewImage(handleWidth, handleHeight)
	handle.Fill(color.RGBA{200, 200, 200, 255})
	handleOpts := &ebiten.DrawImageOptions{}
	handleOpts.GeoM.Translate(float64(handleX), float64(handleY))
	screen.DrawImage(handle, handleOpts)

	// Speed label
	speedLabel := fmt.Sprintf("Speed: %d", 50-g.StepDelay+10)
	// Position speed label to the left of the slider
	text.Draw(screen, speedLabel, basicfont.Face7x13, sliderBgX-text.BoundString(basicfont.Face7x13, speedLabel).Dx()-10, sliderBgY+text.BoundString(basicfont.Face7x13, speedLabel).Dy()/2+basicfont.Face7x13.Ascent/2, color.Black)

	// Draw Zoom level
	zoomLabel := fmt.Sprintf("Zoom: %.1fx", g.ZoomLevel)
	// Position zoom level near the bottom left
	text.Draw(screen, zoomLabel, basicfont.Face7x13, 20, screenHeight-20, color.Black)

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
		text.Draw(screen, title, basicfont.Face7x13, modalX+10, modalY+20, color.Black)

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
		text.Draw(screen, g.AVLInputText, basicfont.Face7x13, inputX+5, inputY+inputHeight/2+basicfont.Face7x13.Ascent/2, color.Black)

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
	bounds := text.BoundString(face, textLabel)
	textX := x + (width-bounds.Dx())/2
	textY := y + (height-bounds.Dy())/2 + basicfont.Face7x13.Ascent
	text.Draw(screen, textLabel, face, textX, textY, textColor)
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
  Mouse Wheel: Zoom In/Out
  Middle Click / Shift+Right Click: Pan
  0: Reset View (Zoom and Pan)
  H: Toggle Help

Context Menu:
  Right Click on node or empty area
`
	text.Draw(screen, helpText, basicfont.Face7x13, helpBgX+20, helpBgY+20, color.Black)

	// Close instruction
	closeText := "Press H to close"
	closeBounds := text.BoundString(basicfont.Face7x13, closeText)
	text.Draw(screen, closeText, basicfont.Face7x13, helpBgX+helpBgWidth-closeBounds.Dx()-20, helpBgY+helpBgHeight-20, color.Black)
}
