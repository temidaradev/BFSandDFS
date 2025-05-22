package ui

import (
	"bfsdfs/internal/algorithms"
	"bfsdfs/internal/graph"
	"bfsdfs/internal/simulator"
	"bfsdfs/pkg/draw"
	"fmt"
	"math"
	"strconv"

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

			// Calculate new offset
			newOffsetX := g.CanvasOffsetX + float64(deltaX)
			newOffsetY := g.CanvasOffsetY + float64(deltaY)

			// Calculate grid boundaries with zoom
			gridSize := float64(1000) // Same as in drawer.go
			minOffset := -gridSize*g.ZoomLevel + float64(screenWidth)
			maxOffset := float64(0)

			// Limit movement within grid boundaries
			if newOffsetX > maxOffset {
				newOffsetX = maxOffset
			} else if newOffsetX < minOffset {
				newOffsetX = minOffset
			}

			if newOffsetY > maxOffset {
				newOffsetY = maxOffset
			} else if newOffsetY < minOffset {
				newOffsetY = minOffset
			}

			// Update offsets
			g.CanvasOffsetX = newOffsetX
			g.CanvasOffsetY = newOffsetY
			g.CanvasDragStartX = g.MouseX
			g.CanvasDragStartY = g.MouseY
			g.canvasNeedsRedraw = true
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

	// Handle AVL Input Modal
	if g.ShowAVLInput {
		// Handle text input (only digits)
		g.AVLInputText += string(ebiten.InputChars())

		// Handle backspace
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(g.AVLInputText) > 0 {
				g.AVLInputText = g.AVLInputText[:len(g.AVLInputText)-1]
			}
		}

		// Handle Enter key (OK)
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if g.AVLAction != "" && g.AVLInputText != "" {
				value, err := strconv.Atoi(g.AVLInputText)
				if err == nil {
					// Perform AVL action
					sim := g.Sim // Use a local variable for clarity
					if sim.GetMode() != algorithms.ModeAVL {
						sim.StartAVL() // Start AVL mode if not already in it
					}
					sim.SetAVLAction(g.AVLAction)
					sim.SetAVLValue(value) // Use the public method
					sim.UpdateAVL()        // Update the AVL tree visualization

					// Perform the actual AVL operation
					switch g.AVLAction {
					case "insert":
						sim.InsertAVL(value)
						g.showMessage(fmt.Sprintf("Inserted %d into AVL tree", value))
					case "delete":
						sim.DeleteAVL(value)
						g.showMessage(fmt.Sprintf("Deleted %d from AVL tree", value))
					case "search":
						sim.SearchAVL(value)
						g.showMessage(fmt.Sprintf("Searched for %d in AVL tree", value))
					}
					sim.UpdateAVL() // Update visualization after operation

					// Close modal
					g.ShowAVLInput = false
				} else {
					g.showMessage("Invalid number")
					g.AVLInputText = "" // Clear invalid input
				}
			} else {
				g.showMessage("Please enter a value")
			}
			return nil // Prevent other actions
		}

		// Handle Escape key (Cancel)
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.ShowAVLInput = false
			return nil // Prevent other actions
		}

		// Handle mouse clicks on OK/Cancel buttons
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			modalWidth := 300
			modalHeight := 150
			modalX := (screenWidth - modalWidth) / 2
			modalY := (screenHeight - modalHeight) / 2
			buttonWidth := 80
			buttonHeight := 30
			buttonSpacing := 10
			buttonY := modalY + modalHeight - buttonHeight - 10

			// Check OK button click
			okButtonX := modalX + modalWidth - buttonWidth*2 - buttonSpacing*2
			if g.MouseX >= okButtonX && g.MouseX <= okButtonX+buttonWidth &&
				g.MouseY >= buttonY && g.MouseY <= buttonY+buttonHeight {
				// Simulate Enter key press
				return g.Update() // Re-run update to process the simulated Enter
			}

			// Check Cancel button click
			cancelButtonX := modalX + modalWidth - buttonWidth - buttonSpacing
			if g.MouseX >= cancelButtonX && g.MouseX <= cancelButtonX+buttonWidth &&
				g.MouseY >= buttonY && g.MouseY <= buttonY+buttonHeight {
				// Simulate Escape key press
				return g.Update() // Re-run update to process the simulated Escape
			}
		}

		return nil // Consume input while modal is open
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

			// If not interacting with buttons, slider, or dialogs, check for node/canvas interaction
			if !g.SliderDragging && !g.ContextMenu.Visible && !g.ShowSaveDialog && !g.ShowLoadDialog && g.MouseY < screenHeight-100 {
				// Convert mouse position to canvas coordinates
				canvasX := float64(g.MouseX) - g.CanvasOffsetX
				canvasY := float64(g.MouseY) - g.CanvasOffsetY

				// Check if clicked on a node for dragging or selection
				targetNode := -1
				for i, node := range g.Sim.Graph.Nodes {
					nodeCanvasX := float64(node.X) * g.ZoomLevel
					nodeCanvasY := float64(node.Y) * g.ZoomLevel
					dx := canvasX - nodeCanvasX
					dy := canvasY - nodeCanvasY
					dist := dx*dx + dy*dy
					if dist <= (20*g.ZoomLevel)*(20*g.ZoomLevel) { // Within the zoomed node radius
						targetNode = i
						break
					}
				}

				if g.EditMode && targetNode != -1 {
					// If in edit mode and clicked on a node, start dragging that node
					g.DraggingNode = targetNode
				} else if targetNode == -1 && !g.EditMode && g.Sim.Mode == algorithms.ModeIdle {
					// If clicked on empty area (not in edit mode and idle), start selection
					g.Selecting = true
					g.SelectionStartX = g.MouseX
					g.SelectionStartY = g.MouseY
					// Clear previous selection if not holding Shift
					if !ebiten.IsKeyPressed(ebiten.KeyShift) {
						g.SelectedNodes = []int{}
						g.SelectedEdges = [][2]int{}
					}
				} else if targetNode != -1 && (isInNodeSelection(g.SelectedNodes, targetNode) || anyEdgeConnectedToNodeIsSelected(g.Sim.Graph, g.SelectedEdges, targetNode)) {
					// If clicked on a selected node or a node connected to a selected edge, start dragging the selection
					g.DraggingSelection = true
					g.SelectionDragStartX = float64(g.MouseX)
					g.SelectionDragStartY = float64(g.MouseY)
				} else if targetNode == -1 && !g.Selecting && !g.DraggingSelection && !ebiten.IsKeyPressed(ebiten.KeyShift) {
					// If clicked on empty area and not selecting/dragging selection and no shift, clear selection
					g.SelectedNodes = []int{}
					g.SelectedEdges = [][2]int{}
				}

				// If clicked on a node and not in edit mode and idle, set it as start node
				if targetNode != -1 && !g.EditMode && g.Sim.Mode == algorithms.ModeIdle {
					g.StartNode = targetNode
					g.showMessage("Start node set to " + string(rune('A'+targetNode)))
				}

				// Handle adding/removing nodes/edges in edit mode
				if g.EditMode {
					if g.RemovingNode {
						if targetNode != -1 {
							// Don't allow removing the last node
							if len(g.Sim.Graph.Nodes) > 1 {
								g.removeNode(targetNode)
								g.showMessage("Node removed")
							} else {
								g.showMessage("Cannot remove the last node")
							}
							g.RemovingNode = false
							g.EditMode = false // Exit edit mode after action
						}
					} else if g.AddingEdge {
						if targetNode != -1 {
							if g.EdgeStartNode == -1 {
								g.EdgeStartNode = targetNode
								g.showMessage("Select second node for the edge")
							} else if g.EdgeStartNode != targetNode {
								g.addEdge(g.EdgeStartNode, targetNode)
								g.EdgeStartNode = -1
								g.AddingEdge = false
								g.EditMode = false // Exit edit mode after action
								g.showMessage("Edge added")
							}
						}
					}
				} else if g.RemovingEdge {
					if targetNode != -1 {
						if g.EdgeStartNode == -1 {
							g.EdgeStartNode = targetNode
							g.showMessage("Select second node to remove edge")
						} else if g.EdgeStartNode != targetNode {
							g.removeEdge(g.EdgeStartNode, targetNode)
							g.EdgeStartNode = -1
							g.RemovingEdge = false
							g.EditMode = false // Exit edit mode after action
						}
					}
				} else if targetNode == -1 && len(g.Sim.Graph.Nodes) < 15 && !g.Selecting && !g.DraggingSelection {
					// If clicked on empty area and not selecting/dragging selection, add node
					// Snap to grid if enabled (in canvas coordinates)
					nodeX, nodeY := int(canvasX/g.ZoomLevel), int(canvasY/g.ZoomLevel)
					if g.SnapToGrid {
						nodeX, nodeY = draw.SnapToGrid(nodeX, nodeY, g.GridConfig.CellSize)
					}

					// Keep node within grid boundaries (adjusting for zoom)
					gridSize := 1000.0 // Same as in drawer.go
					minCanvasX := 20.0 / g.ZoomLevel
					maxCanvasX := (gridSize - 20.0) / g.ZoomLevel
					minCanvasY := 20.0 / g.ZoomLevel
					maxCanvasY := (gridSize - 20.0) / g.ZoomLevel

					nodeX = int(math.Max(minCanvasX, math.Min(maxCanvasX, float64(nodeX))))
					nodeY = int(math.Max(minCanvasY, math.Min(maxCanvasY, float64(nodeY))))

					g.addNode(nodeX, nodeY)
					g.showMessage("Node added")
				}
			}
		}
		g.MouseClicked = true // Set to true as mouse button is pressed

	} else {
		// Mouse released
		if g.MouseClicked {
			// If released after selecting, finalize selection
			if g.Selecting {
				g.finalizeSelection(g.SelectionStartX, g.SelectionStartY, g.MouseX, g.MouseY)
				g.Selecting = false
			}
			g.MouseReleased = true // Set to true as mouse button is released
		}
		g.MouseClicked = false
		g.SliderDragging = false
		g.DraggingNode = -1
		g.DraggingSelection = false // Stop dragging selection on mouse release
		g.MouseReleased = false
	}

	// Handle dragging a node (if not dragging a selection)
	if g.DraggingNode != -1 && !g.DraggingSelection {
		// Convert mouse position to canvas coordinates
		canvasX := float64(g.MouseX) - g.CanvasOffsetX
		canvasY := float64(g.MouseY) - g.CanvasOffsetY

		// Convert canvas coordinates back to node coordinates (undoing zoom)
		nodeX := int(canvasX / g.ZoomLevel)
		nodeY := int(canvasY / g.ZoomLevel)

		// Snap to grid if enabled
		if g.SnapToGrid {
			nodeX, nodeY = draw.SnapToGrid(nodeX, nodeY, g.GridConfig.CellSize)
		}

		// Keep node within grid boundaries
		gridSize := 1000.0 // Same as in drawer.go
		minCanvasX := 20.0 / g.ZoomLevel
		maxCanvasX := (gridSize - 20.0) / g.ZoomLevel
		minCanvasY := 20.0 / g.ZoomLevel
		maxCanvasY := (gridSize - 20.0) / g.ZoomLevel

		// Convert grid boundaries back to node coordinate system for clamping
		minNodeX := int(minCanvasX)
		maxNodeX := int(maxCanvasX)
		minNodeY := int(minCanvasY)
		maxNodeY := int(maxCanvasY)

		nodeX = int(math.Max(float64(minNodeX), math.Min(float64(maxNodeX), float64(nodeX))))
		nodeY = int(math.Max(float64(minNodeY), math.Min(float64(maxNodeY), float64(nodeY))))

		g.Sim.Graph.Nodes[g.DraggingNode].X = nodeX
		g.Sim.Graph.Nodes[g.DraggingNode].Y = nodeY
		g.canvasNeedsRedraw = true
	}

	// Handle dragging a selection
	if g.DraggingSelection {
		// Calculate movement delta in screen coordinates
		deltaX := float64(g.MouseX) - g.SelectionDragStartX
		deltaY := float64(g.MouseY) - g.SelectionDragStartY

		// Move all selected nodes
		for _, nodeIndex := range g.SelectedNodes {
			// Convert current node position to screen coordinates
			nodeScreenX := float64(g.Sim.Graph.Nodes[nodeIndex].X)*g.ZoomLevel + g.CanvasOffsetX
			nodeScreenY := float64(g.Sim.Graph.Nodes[nodeIndex].Y)*g.ZoomLevel + g.CanvasOffsetY

			// Calculate new screen position after applying delta
			newNodeScreenX := nodeScreenX + deltaX
			newNodeScreenY := nodeScreenY + deltaY

			// Convert new screen position back to node coordinates (undoing offset and zoom)
			newNodeX := int((newNodeScreenX - g.CanvasOffsetX) / g.ZoomLevel)
			newNodeY := int((newNodeScreenY - g.CanvasOffsetY) / g.ZoomLevel)

			// Snap to grid if enabled (apply to the new node coordinates)
			if g.SnapToGrid {
				newNodeX, newNodeY = draw.SnapToGrid(newNodeX, newNodeY, g.GridConfig.CellSize)
			}

			// Keep node within grid boundaries (adjusting for zoom is already done in calculation)
			gridSize := 1000.0 // Same as in drawer.go
			minNodeCoord := 20.0
			maxNodeCoord := gridSize - 20.0

			newNodeX = int(math.Max(minNodeCoord, math.Min(maxNodeCoord, float64(newNodeX))))
			newNodeY = int(math.Max(minNodeCoord, math.Min(maxNodeCoord, float64(newNodeY))))

			// Update node position
			g.Sim.Graph.Nodes[nodeIndex].X = newNodeX
			g.Sim.Graph.Nodes[nodeIndex].Y = newNodeY
		}

		// Update drag start position for the next frame
		g.SelectionDragStartX = float64(g.MouseX)
		g.SelectionDragStartY = float64(g.MouseY)
		g.canvasNeedsRedraw = true
	}

	// Auto-stepping
	if g.AutoStep && !g.Sim.Done && g.Sim.Mode != algorithms.ModeIdle && g.Sim.Mode != algorithms.ModeAVL {
		g.StepCounter++
		if g.StepCounter >= g.StepDelay {
			g.StepCounter = 0
			g.Sim.Update()
		}
	} else if g.AutoStep && g.Sim.Mode == algorithms.ModeAVL {
		// Disable auto-stepping when in AVL mode
		g.AutoStep = false
	}

	// Handle keyboard controls for convenience
	handleKeyboardInput(g)

	// Handle help toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.ShowHelp = !g.ShowHelp
	}

	// Handle zoom with mouse wheel
	if _, wheelY := ebiten.Wheel(); wheelY != 0 {
		// Get mouse position in canvas coordinates
		canvasX := g.MouseX - int(g.CanvasOffsetX)
		canvasY := g.MouseY - int(g.CanvasOffsetY)

		// Calculate zoom factor (0.1 per wheel step)
		zoomFactor := 1.0 + wheelY*0.1
		newZoom := g.ZoomLevel * zoomFactor

		// Limit zoom range (0.5x to 2.0x)
		if newZoom >= 0.5 && newZoom <= 2.0 {
			// Calculate new offset to zoom towards mouse position
			newOffsetX := float64(g.MouseX) - float64(canvasX)*newZoom
			newOffsetY := float64(g.MouseY) - float64(canvasY)*newZoom

			// Get screen dimensions
			screenWidth, _ := ebiten.WindowSize()

			// Calculate grid boundaries
			gridSize := float64(1000) // Same as in drawer.go
			minOffset := -gridSize*newZoom + float64(screenWidth)
			maxOffset := float64(0)

			// Ensure new offset stays within grid boundaries
			if newOffsetX > maxOffset {
				newOffsetX = maxOffset
			} else if newOffsetX < minOffset {
				newOffsetX = minOffset
			}

			if newOffsetY > maxOffset {
				newOffsetY = maxOffset
			} else if newOffsetY < minOffset {
				newOffsetY = minOffset
			}

			// Update zoom and offset
			g.ZoomLevel = newZoom
			g.CanvasOffsetX = newOffsetX
			g.CanvasOffsetY = newOffsetY
		}
	}

	// Handle zoom with keyboard shortcuts
	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		// Zoom in
		newZoom := g.ZoomLevel * 1.1
		if newZoom <= 2.0 {
			// Get screen dimensions
			screenWidth, _ := ebiten.WindowSize()

			// Calculate grid boundaries
			gridSize := float64(1000) // Same as in drawer.go
			minOffset := -gridSize*newZoom + float64(screenWidth)
			maxOffset := float64(0)

			// Calculate new offset to keep center point
			centerX := float64(screenWidth/2) - g.CanvasOffsetX
			centerY := float64(screenWidth/2) - g.CanvasOffsetY
			newOffsetX := float64(screenWidth/2) - centerX*newZoom
			newOffsetY := float64(screenWidth/2) - centerY*newZoom

			// Ensure new offset stays within grid boundaries
			if newOffsetX > maxOffset {
				newOffsetX = maxOffset
			} else if newOffsetX < minOffset {
				newOffsetX = minOffset
			}

			if newOffsetY > maxOffset {
				newOffsetY = maxOffset
			} else if newOffsetY < minOffset {
				newOffsetY = minOffset
			}

			// Update zoom and offset
			g.ZoomLevel = newZoom
			g.CanvasOffsetX = newOffsetX
			g.CanvasOffsetY = newOffsetY
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) {
		// Zoom out
		newZoom := g.ZoomLevel * 0.9
		if newZoom >= 0.5 {
			// Get screen dimensions
			screenWidth, _ := ebiten.WindowSize()

			// Calculate grid boundaries
			gridSize := float64(1000) // Same as in drawer.go
			minOffset := -gridSize*newZoom + float64(screenWidth)
			maxOffset := float64(0)

			// Calculate new offset to keep center point
			centerX := float64(screenWidth/2) - g.CanvasOffsetX
			centerY := float64(screenWidth/2) - g.CanvasOffsetY
			newOffsetX := float64(screenWidth/2) - centerX*newZoom
			newOffsetY := float64(screenWidth/2) - centerY*newZoom

			// Ensure new offset stays within grid boundaries
			if newOffsetX > maxOffset {
				newOffsetX = maxOffset
			} else if newOffsetX < minOffset {
				newOffsetX = minOffset
			}

			if newOffsetY > maxOffset {
				newOffsetY = maxOffset
			} else if newOffsetY < minOffset {
				newOffsetY = minOffset
			}

			// Update zoom and offset
			g.ZoomLevel = newZoom
			g.CanvasOffsetX = newOffsetX
			g.CanvasOffsetY = newOffsetY
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key0) {
		// Reset zoom and center the view
		screenWidth, screenHeight := ebiten.WindowSize()
		gridSize := float64(1000) // Same as in drawer.go
		g.ZoomLevel = 1.0
		g.CanvasOffsetX = (float64(screenWidth) - gridSize) / 2
		g.CanvasOffsetY = (float64(screenHeight) - gridSize) / 2
	}

	return nil
}

// Helper function to check if a node is in the selected nodes list
func isInNodeSelection(selectedNodes []int, nodeIndex int) bool {
	for _, index := range selectedNodes {
		if index == nodeIndex {
			return true
		}
	}
	return false
}

// Helper function to get edges connected to a node
func getEdgesConnectedToNode(graph graph.Graph, nodeIndex int) [][2]int {
	var connectedEdges [][2]int
	for _, edge := range graph.Edges {
		if edge[0] == nodeIndex || edge[1] == nodeIndex {
			connectedEdges = append(connectedEdges, edge)
		}
	}
	return connectedEdges
}

// Helper function to check if an edge is in the selected edges list
func isInEdgeSelection(selectedEdges [][2]int, edge [2]int) bool {
	for _, selectedEdge := range selectedEdges {
		// Check for both directions of the edge
		if (selectedEdge[0] == edge[0] && selectedEdge[1] == edge[1]) || (selectedEdge[0] == edge[1] && selectedEdge[1] == edge[0]) {
			return true
		}
	}
	return false
}

// Helper function to check if any edge connected to a node is selected
func anyEdgeConnectedToNodeIsSelected(graph graph.Graph, selectedEdges [][2]int, nodeIndex int) bool {
	// Get all edges connected to the node
	connectedEdges := getEdgesConnectedToNode(graph, nodeIndex)

	// Check if any of these connected edges are in the selected edges list
	for _, edge := range connectedEdges {
		if isInEdgeSelection(selectedEdges, edge) {
			return true
		}
	}

	return false
}

// finalizeSelection determines which nodes and edges are within the selection box
func (g *Game) finalizeSelection(startX, startY, endX, endY int) {
	// Determine the boundaries of the selection box in screen coordinates
	left := min(startX, endX)
	right := max(startX, endX)
	top := min(startY, endY)
	bottom := max(startY, endY)

	// Clear previous selection if Shift key is not held
	if !ebiten.IsKeyPressed(ebiten.KeyShift) {
		g.SelectedNodes = []int{}
		g.SelectedEdges = [][2]int{}
	}

	// Identify nodes within the selection box
	for i, node := range g.Sim.Graph.Nodes {
		// Convert node position to screen coordinates
		nodeScreenX := int(float64(node.X)*g.ZoomLevel + g.CanvasOffsetX)
		nodeScreenY := int(float64(node.Y)*g.ZoomLevel + g.CanvasOffsetY)

		// Check if node is within the selection box boundaries
		if nodeScreenX >= left && nodeScreenX <= right && nodeScreenY >= top && nodeScreenY <= bottom {
			// Add node to selection if not already selected
			if !isInNodeSelection(g.SelectedNodes, i) {
				g.SelectedNodes = append(g.SelectedNodes, i)
			}
		}
	}

	// Identify edges within the selection box
	for _, edge := range g.Sim.Graph.Edges {
		// Get the connected nodes
		node1 := g.Sim.Graph.Nodes[edge[0]]
		node2 := g.Sim.Graph.Nodes[edge[1]]

		// Convert node positions to screen coordinates
		x1 := float64(node1.X)*g.ZoomLevel + g.CanvasOffsetX
		y1 := float64(node1.Y)*g.ZoomLevel + g.CanvasOffsetY
		x2 := float64(node2.X)*g.ZoomLevel + g.CanvasOffsetX
		y2 := float64(node2.Y)*g.ZoomLevel + g.CanvasOffsetY

		// Check if the edge intersects the selection box
		// A simple check: if both endpoints are within the box, select the edge.
		// More complex line-box intersection could be added later if needed.
		isEndpoint1InBox := x1 >= float64(left) && x1 <= float64(right) && y1 >= float64(top) && y1 <= float64(bottom)
		isEndpoint2InBox := x2 >= float64(left) && x2 <= float64(right) && y2 >= float64(top) && y2 <= float64(bottom)

		if isEndpoint1InBox || isEndpoint2InBox {
			// Add edge to selection if not already selected
			if !isInEdgeSelection(g.SelectedEdges, edge) {
				g.SelectedEdges = append(g.SelectedEdges, edge)
			}
		}
	}

	g.canvasNeedsRedraw = true
}
