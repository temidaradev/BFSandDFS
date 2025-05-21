package ui

import (
	//"fmt"
	"image/color"

	"bfsdfs/internal/algorithms"
	"bfsdfs/pkg/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Get window dimensions
	screenWidth, screenHeight := ebiten.WindowSize()

	// Create a separate canvas for graph visualization - use full screen instead of reserving space
	graphCanvas := ebiten.NewImage(screenWidth, screenHeight)
	graphCanvas.Fill(color.RGBA{240, 240, 240, 255})

	// Draw grid if enabled
	if g.ShowGrid {
		draw.DrawGrid(graphCanvas, screenWidth, screenHeight, g.GridConfig)
	}

	// Draw edges
	for _, edge := range g.Sim.Graph.Edges {
		// Get the connected nodes
		node1 := g.Sim.Graph.Nodes[edge[0]]
		node2 := g.Sim.Graph.Nodes[edge[1]]

		// Determine edge color based on visited state
		edgeColor := color.RGBA{150, 150, 150, 255} // Default gray

		// Check if both nodes are visited
		if g.Sim.Mode != algorithms.ModeIdle {
			if g.Sim.Visited[edge[0]] && g.Sim.Visited[edge[1]] {
				edgeColor = color.RGBA{100, 180, 100, 255} // Green for visited
			}
		}

		// Apply canvas offset for movement
		x1 := float64(node1.X) + g.CanvasOffsetX
		y1 := float64(node1.Y) + g.CanvasOffsetY
		x2 := float64(node2.X) + g.CanvasOffsetX
		y2 := float64(node2.Y) + g.CanvasOffsetY

		// Draw the edge
		draw.DrawLine(graphCanvas, x1, y1, x2, y2, edgeColor)
	}

	// Draw nodes
	for i, node := range g.Sim.Graph.Nodes {
		// Determine node color based on state
		var nodeColor color.RGBA

		if g.StartNode == i {
			nodeColor = color.RGBA{60, 120, 200, 255} // Blue for start node
		} else if g.Sim.Mode != algorithms.ModeIdle && g.Sim.Current == i {
			nodeColor = color.RGBA{220, 100, 100, 255} // Red for current
		} else if g.Sim.Mode != algorithms.ModeIdle && g.Sim.LastActive == i {
			nodeColor = color.RGBA{200, 150, 100, 255} // Orange for last active
		} else if g.Sim.Mode != algorithms.ModeIdle && g.Sim.Visited[i] {
			nodeColor = color.RGBA{100, 180, 100, 255} // Green for visited
		} else {
			nodeColor = color.RGBA{200, 200, 200, 255} // Light gray for unvisited
		}

		// Apply canvas offset for movement
		nodeX := node.X + int(g.CanvasOffsetX)
		nodeY := node.Y + int(g.CanvasOffsetY)

		// Draw the node
		draw.DrawCircle(graphCanvas, nodeX, nodeY, 20, nodeColor)

		// Draw the node label (A, B, C, etc.)
		label := string(rune('A' + i))
		text.Draw(graphCanvas, label, basicfont.Face7x13, nodeX-4, nodeY+4, color.Black)
	}

	// Draw the graph canvas onto the main screen
	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(graphCanvas, opts)

	// Draw hint text about canvas movement - use simple ASCII arrow for compatibility
	hintText := "Middle-click or Shift+Right-click and drag to move the canvas"
	text.Draw(screen, hintText, basicfont.Face7x13, screenWidth/2-160, 20, color.RGBA{30, 30, 150, 255})

	// Draw algorithm info if active
	if g.Sim.Mode != algorithms.ModeIdle {
		// Draw visit order in a semi-transparent overlay at the top
		orderStr := "Visit order: "
		for i, nodeIdx := range g.Sim.Order {
			if i > 0 {
				orderStr += " > " // Changed from → to > for compatibility
			}
			orderStr += string(rune('A' + nodeIdx))
		}

		// Create a dark background for the text
		textBg := ebiten.NewImage(screenWidth, 20)
		textBg.Fill(color.RGBA{40, 40, 40, 200})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, 30)
		screen.DrawImage(textBg, opts)

		text.Draw(screen, orderStr, basicfont.Face7x13, 25, 45, color.White)

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

		// Create a background for this text too
		dsBg := ebiten.NewImage(screenWidth, 20)
		dsBg.Fill(color.RGBA{40, 40, 40, 200})
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, 55)
		screen.DrawImage(dsBg, opts)

		text.Draw(screen, dataStructStr, basicfont.Face7x13, 25, 70, color.White)
	}

	// Draw all buttons directly on the screen
	for _, btn := range g.Buttons {
		btn.Draw(screen, g)
	}

	// Draw speed slider if algorithm is running
	if g.Sim.Mode != algorithms.ModeIdle {
		sliderY := screenHeight - 50
		sliderBg := ebiten.NewImage(200, 20)
		sliderBg.Fill(color.RGBA{40, 40, 40, 220})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(screenWidth-220), float64(sliderY))
		screen.DrawImage(sliderBg, opts)

		sliderPos := float64(screenWidth-220) + float64(200-((g.StepDelay-10)*200/40))
		sliderHandle := ebiten.NewImage(10, 20)
		sliderHandle.Fill(color.RGBA{100, 100, 200, 255})
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(sliderPos, float64(sliderY))
		screen.DrawImage(sliderHandle, opts)

		// Draw a small label with speed
		speedBg := ebiten.NewImage(50, 20)
		speedBg.Fill(color.RGBA{40, 40, 40, 220})
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(screenWidth-270), float64(sliderY))
		screen.DrawImage(speedBg, opts)

		text.Draw(screen, "Speed:", basicfont.Face7x13, screenWidth-265, sliderY+15, color.White)
	}

	// Draw temporary message with a semi-transparent background
	if g.MessageTimer > 0 {
		msgWidth := len(g.Message) * 7 // Approximation
		msgBg := ebiten.NewImage(msgWidth+20, 20)
		msgBg.Fill(color.RGBA{40, 40, 40, 200})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(20, float64(screenHeight-85))
		screen.DrawImage(msgBg, opts)

		messageColor := color.RGBA{220, 220, 100, 255} // Yellow for better visibility
		text.Draw(screen, g.Message, basicfont.Face7x13, 30, screenHeight-70, messageColor)
	}

	// Draw dialogs
	if g.ShowSaveDialog {
		g.SaveDialog.Draw(screen)
	}

	if g.ShowLoadDialog {
		g.LoadDialog.Draw(screen)
	}

	// Draw context menu
	g.ContextMenu.Draw(screen)
}
