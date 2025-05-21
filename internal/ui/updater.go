package ui

import (
	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/graph"
	"bfsdfs/internal/simulator"
	"bfsdfs/pkg/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Update handles user input and simulation updates
func (g *Game) Update() error {
	// Get window dimensions for calculations
	screenWidth, screenHeight := ebiten.WindowSize()

	// Track mouse position
	g.MouseX, g.MouseY = ebiten.CursorPosition()

	// Update message timer
	if g.MessageTimer > 0 {
		g.MessageTimer--
	}

	// Handle button hover state
	for _, btn := range g.Buttons {
		// Calculate button position based on anchoring
		btnX, btnY := g.getAdjustedButtonPosition(btn)
		btn.Hover = g.MouseX >= btnX && g.MouseX <= btnX+btn.Width &&
			g.MouseY >= btnY && g.MouseY <= btnY+btn.Height
	}

	// Update context menu hover states
	g.ContextMenu.UpdateHoverState(g.MouseX, g.MouseY)

	// Handle canvas dragging (middle mouse button or right mouse button with shift)
	if (inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) ||
		(inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && ebiten.IsKeyPressed(ebiten.KeyShift))) &&
		g.MouseY < screenHeight-100 {
		g.CanvasDragging = true
		g.CanvasDragStartX = g.MouseX
		g.CanvasDragStartY = g.MouseY
		return nil // Prevent other actions when starting drag
	}

	if g.CanvasDragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) ||
			(ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ebiten.IsKeyPressed(ebiten.KeyShift)) {
			// Update canvas offset based on mouse movement
			deltaX := g.MouseX - g.CanvasDragStartX
			deltaY := g.MouseY - g.CanvasDragStartY
			g.CanvasOffsetX += float64(deltaX)
			g.CanvasOffsetY += float64(deltaY)
			g.CanvasDragStartX = g.MouseX
			g.CanvasDragStartY = g.MouseY
			return nil // Prevent other actions while dragging
		} else {
			g.CanvasDragging = false
		}
	}

	// Handle right-click for context menu
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && !ebiten.IsKeyPressed(ebiten.KeyShift) && g.MouseY < screenHeight-100 {
		// Calculate mouse position in canvas coordinates (accounting for offset)
		canvasX := g.MouseX - int(g.CanvasOffsetX)
		canvasY := g.MouseY - int(g.CanvasOffsetY)

		// Check if right-clicked on a node
		targetNode := -1
		for i, node := range g.Sim.Graph.Nodes {
			dx := float64(canvasX - node.X)
			dy := float64(canvasY - node.Y)
			dist := dx*dx + dy*dy
			if dist <= 20*20 { // Within node radius
				targetNode = i
				break
			}
		}

		// Show context menu with appropriate options
		g.ContextMenu.ClearItems()

		if targetNode != -1 {
			// Node-specific options
			g.ContextMenu.AddItem("Set as Start Node", func() {
				g.StartNode = targetNode
				g.showMessage("Start node set to " + string(rune('A'+targetNode)))
			})

			g.ContextMenu.AddItem("Delete Node", func() {
				// Don't allow removing the last node
				if len(g.Sim.Graph.Nodes) > 1 {
					g.removeNode(targetNode)
					g.showMessage("Node removed")
				} else {
					g.showMessage("Cannot remove the last node")
				}
			})

			g.ContextMenu.AddItem("Add Edge From Here", func() {
				g.EdgeStartNode = targetNode
				g.AddingEdge = true
				g.EditMode = true
				g.showMessage("Select second node for the edge")
			})

			// Add an option to remove all edges from this node
			g.ContextMenu.AddItem("Clear Node Edges", func() {
				g.clearNodeEdges(targetNode)
				g.showMessage("Cleared all edges from node " + string(rune('A'+targetNode)))
			})
		} else {
			// Empty area options
			g.ContextMenu.AddItem("Add Node Here", func() {
				if len(g.Sim.Graph.Nodes) < 15 {
					// Get canvas coordinates and snap to grid if needed
					nodeX, nodeY := canvasX, canvasY
					if g.SnapToGrid {
						nodeX, nodeY = draw.SnapToGrid(nodeX, nodeY, g.GridConfig.CellSize)
					}
					g.addNode(nodeX, nodeY)
					g.showMessage("Node added")
				} else {
					g.showMessage("Maximum node count reached (15)")
				}
			})

			g.ContextMenu.AddItem("Create Random Graph", func() {
				nodeCount := len(g.Sim.Graph.Nodes)
				*g.Sim = *simulator.NewSimulator(nodeCount)
				g.StartNode = 0
				g.showMessage("Random graph created")
			})
		}

		// Add save/load options
		g.ContextMenu.AddItem("Save Graph...", func() {
			g.SaveDialog.Show()
			g.ShowSaveDialog = true
		})

		g.ContextMenu.AddItem("Load Graph...", func() {
			g.LoadDialog.Show()
			g.ShowLoadDialog = true
		})

		// Add general options
		g.ContextMenu.AddItem("Clear All Edges", func() {
			// Clear all edges but keep nodes
			g.Sim.Graph.Edges = [][2]int{}
			for i := range g.Sim.Graph.Nodes {
				g.Sim.Graph.Nodes[i].Neighbors = []int{}
			}
			g.showMessage("All edges cleared")
		})

		g.ContextMenu.Show(g.MouseX, g.MouseY, targetNode)
		return nil
	}

	// Handle save dialog
	if g.ShowSaveDialog {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Handle OK button click in save dialog
			if g.MouseX >= g.SaveDialog.X+g.SaveDialog.Width-180 &&
				g.MouseX <= g.SaveDialog.X+g.SaveDialog.Width-100 &&
				g.MouseY >= g.SaveDialog.Y+g.SaveDialog.Height-30 &&
				g.MouseY <= g.SaveDialog.Y+g.SaveDialog.Height {

				// Save the graph to the selected file
				filePath := g.SaveDialog.GetSelectedFilePath()
				if err := g.Sim.Graph.SaveGraph(filePath); err != nil {
					g.showMessage("Error saving graph: " + err.Error())
				} else {
					g.showMessage("Graph saved to " + filePath)
				}
				g.SaveDialog.Hide()
				g.ShowSaveDialog = false
				return nil
			}

			// Let the file dialog handle other clicks
			if g.SaveDialog.HandleClick(g.MouseX, g.MouseY) {
				return nil
			}

			// Close dialog if clicked outside
			g.SaveDialog.Hide()
			g.ShowSaveDialog = false
		}

		// Handle keyboard input for save dialog
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.SaveDialog.Hide()
			g.ShowSaveDialog = false
			return nil
		}

		return nil
	}

	// Handle load dialog
	if g.ShowLoadDialog {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Handle OK button click in load dialog
			if g.MouseX >= g.LoadDialog.X+g.LoadDialog.Width-180 &&
				g.MouseX <= g.LoadDialog.X+g.LoadDialog.Width-100 &&
				g.MouseY >= g.LoadDialog.Y+g.LoadDialog.Height-30 &&
				g.MouseY <= g.LoadDialog.Y+g.LoadDialog.Height {

				// Load the graph from the selected file
				filePath := g.LoadDialog.GetSelectedFilePath()
				loadedGraph, err := graph.LoadGraph(filePath)
				if err != nil {
					g.showMessage("Error loading graph: " + err.Error())
				} else {
					g.Sim.Graph = *loadedGraph
					g.Sim.Reset()
					g.StartNode = 0
					g.showMessage("Graph loaded from " + filePath)
				}
				g.LoadDialog.Hide()
				g.ShowLoadDialog = false
				return nil
			}

			// Let the file dialog handle other clicks
			if g.LoadDialog.HandleClick(g.MouseX, g.MouseY) {
				return nil
			}

			// Close dialog if clicked outside
			g.LoadDialog.Hide()
			g.ShowLoadDialog = false
		}

		// Handle keyboard input for load dialog
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.LoadDialog.Hide()
			g.ShowLoadDialog = false
			return nil
		}

		return nil
	}

	// Handle context menu clicks
	if g.ContextMenu.Visible && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.ContextMenu.HandleClick(g.MouseX, g.MouseY) {
			return nil
		}
	}

	// Handle left mouse press
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Handle button clicks when mouse is first pressed
		if !g.MouseClicked {
			// Check for button clicks using adjusted positions
			for _, btn := range g.Buttons {
				btnX, btnY := g.getAdjustedButtonPosition(btn)
				if g.MouseX >= btnX && g.MouseX <= btnX+btn.Width &&
					g.MouseY >= btnY && g.MouseY <= btnY+btn.Height {
					btn.Action()
					g.MouseClicked = true
					return nil
				}
			}

			// Check for slider interaction in the HUD area
			if g.MouseY >= screenHeight-60 && g.MouseY <= screenHeight-40 &&
				g.MouseX >= screenWidth-220 && g.MouseX <= screenWidth-20 {
				g.SliderDragging = true
				// Update slider position immediately
				g.StepDelay = 50 - int((float64(g.MouseX-(screenWidth-220))/200.0)*40)
				if g.StepDelay < 10 {
					g.StepDelay = 10
				} else if g.StepDelay > 50 {
					g.StepDelay = 50
				}
			}

			// Handle node interaction in graph area, accounting for canvas offset
			if g.EditMode && g.DraggingNode == -1 && !g.SliderDragging &&
				g.Sim.Mode == algorithms.ModeIdle && !g.ContextMenu.Visible &&
				g.MouseY < screenHeight-100 {

				// Convert mouse position to canvas coordinates
				canvasX := g.MouseX - int(g.CanvasOffsetX)
				canvasY := g.MouseY - int(g.CanvasOffsetY)

				// Check if clicked on a node
				for i, node := range g.Sim.Graph.Nodes {
					dx := float64(canvasX - node.X)
					dy := float64(canvasY - node.Y)
					dist := dx*dx + dy*dy
					if dist <= 20*20 { // Within the node radius
						if g.RemovingNode {
							// Don't allow removing the last node
							if len(g.Sim.Graph.Nodes) > 1 {
								g.removeNode(i)
								g.showMessage("Node removed")
							} else {
								g.showMessage("Cannot remove the last node")
							}
							g.RemovingNode = false
							g.EditMode = false
						} else if g.AddingEdge {
							if g.EdgeStartNode == -1 {
								g.EdgeStartNode = i
								g.showMessage("Select second node for the edge")
							} else if g.EdgeStartNode != i {
								g.addEdge(g.EdgeStartNode, i)
								g.EdgeStartNode = -1
								g.AddingEdge = false
								g.EditMode = false
								g.showMessage("Edge added")
							}
						} else if g.RemovingEdge {
							if g.EdgeStartNode == -1 {
								g.EdgeStartNode = i
								g.showMessage("Select second node to remove edge")
							} else if g.EdgeStartNode != i {
								g.removeEdge(g.EdgeStartNode, i)
								g.EdgeStartNode = -1
								g.RemovingEdge = false
								g.EditMode = false
							}
						} else {
							// Start dragging the node
							g.DraggingNode = i
						}
						break
					}
				}

				// If clicked on empty area in edit mode
				if g.DraggingNode == -1 && !g.RemovingNode && !g.AddingEdge && !g.RemovingEdge &&
					g.MouseY < screenHeight-100 && len(g.Sim.Graph.Nodes) < 15 {
					// Snap to grid if enabled (in canvas coordinates)
					nodeX, nodeY := canvasX, canvasY
					if g.SnapToGrid {
						nodeX, nodeY = draw.SnapToGrid(nodeX, nodeY, g.GridConfig.CellSize)
					}
					g.addNode(nodeX, nodeY)
					g.showMessage("Node added")
				}
			}

			// Set start node when not in edit mode
			if !g.EditMode && g.Sim.Mode == algorithms.ModeIdle && !g.SliderDragging &&
				!g.ContextMenu.Visible && g.MouseY < screenHeight-100 {
				// Convert mouse position to canvas coordinates
				canvasX := g.MouseX - int(g.CanvasOffsetX)
				canvasY := g.MouseY - int(g.CanvasOffsetY)

				for i, node := range g.Sim.Graph.Nodes {
					dx := float64(canvasX - node.X)
					dy := float64(canvasY - node.Y)
					dist := dx*dx + dy*dy
					if dist <= 20*20 { // Within the circle radius
						g.StartNode = i
						g.showMessage("Start node set to " + string(rune('A'+i)))
						break
					}
				}
			}
		}

		// Handle dragging a node
		if g.DraggingNode != -1 {
			// Convert mouse position to canvas coordinates
			canvasX := g.MouseX - int(g.CanvasOffsetX)
			canvasY := g.MouseY - int(g.CanvasOffsetY)

			// Snap to grid if enabled
			if g.SnapToGrid {
				canvasX, canvasY = draw.SnapToGrid(canvasX, canvasY, g.GridConfig.CellSize)
			}

			// Keep node within reasonable bounds
			if canvasY < 20 {
				canvasY = 20
			} else if canvasY > screenHeight-120 {
				canvasY = screenHeight - 120
			}

			g.Sim.Graph.Nodes[g.DraggingNode].X = canvasX
			g.Sim.Graph.Nodes[g.DraggingNode].Y = canvasY
		}

		// Update slider position if dragging
		if g.SliderDragging {
			sliderX := g.MouseX
			if sliderX < screenWidth-220 {
				sliderX = screenWidth - 220
			}
			if sliderX > screenWidth-20 {
				sliderX = screenWidth - 20
			}
			// Map slider position to delay (50 to 10 frames)
			g.StepDelay = 50 - int((float64(sliderX-(screenWidth-220))/200.0)*40)
		}

		g.MouseClicked = true
	} else {
		// Mouse released
		if g.MouseClicked {
			g.MouseReleased = true
		}
		g.MouseClicked = false
		g.SliderDragging = false
		g.DraggingNode = -1
		g.MouseReleased = false
	}

	// Auto-stepping
	if g.AutoStep && !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle {
		g.StepCounter++
		if g.StepCounter >= g.StepDelay {
			g.StepCounter = 0
			g.Sim.Update()
		}
	}

	// Keep keyboard controls for convenience
	handleKeyboardInput(g)

	return nil
}
