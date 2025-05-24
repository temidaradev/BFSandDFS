# BFS, DFS, and AVL Tree Simulator

This application demonstrates Breadth-First Search (BFS), Depth-First Search (DFS) graph traversal algorithms, and AVL tree operations through an interactive visualization.

## Project Structure

The project is organized as follows:

```
.
├── cmd/
│   └── simulator/       # Main application entry point
│       └── main.go
├── internal/
│   ├── algorithms/      # Algorithm implementations
│   │   └── traversal.go
│   ├── graph/           # Graph data structures
│   │   └── graph.go
│   ├── simulator/       # Simulator core logic
│   │   └── simulator.go
│   └── ui/              # User interface components
│       └── game.go
└── pkg/
    └── draw/            # Reusable drawing utilities
        └── helpers.go
```

## How to Run

The easiest way to run the simulator is with the provided Makefile:

```bash
# Run the application
make run

# Build the application
make build

# Build and immediately run
make dev

# Clean build artifacts
make clean
```

Or run directly with Go:

```bash
go run ./cmd/simulator
```

## UI Controls

### Control Buttons

- **BFS**: Start Breadth-First Search from the selected node
- **DFS**: Start Depth-First Search from the selected node
- **AVL Tree**: Switch to AVL tree mode for tree operations
- **Step**: Perform one step of the algorithm (BFS/DFS only)
- **Auto**: Toggle automatic stepping (BFS/DFS only)
- **Reset**: Reset the simulation to initial state

### AVL Tree Operation Buttons (visible in AVL mode)

- **Insert**: Insert a value into the AVL tree (prompts for input)
- **Delete**: Delete a value from the AVL tree (prompts for input)
- **Search**: Search for a value in the AVL tree (prompts for input)

### Graph Editing Buttons

- **New Graph**: Generate a new random graph with the same number of nodes
- **Add Node**: Add one more node to the graph (max 15)
- **Del Node**: Enter node deletion mode - click a node to remove it
- **Add Edge**: Enter edge creation mode - click two nodes to add an edge between them
- **Del Edge**: Enter edge deletion mode - click two nodes to remove the edge between them
- **Edit Mode**: Toggle edit mode where you can drag nodes to reposition them
- **Grid**: Toggle grid display for easier node positioning
- **Snap**: Toggle snap-to-grid feature for precise node placement

### File Operation Buttons

- **Save**: Open the save dialog to save the current graph to a JSON file
- **Load**: Open the load dialog to load a graph from a JSON file

### Context Menu

Right-click on the graph area to access a context menu with additional options:

- **Right-click on a node**:
  - Set as Start Node: Makes the node the starting point for traversals
  - Delete Node: Removes the node from the graph
  - Add Edge From Here: Starts the edge creation process from this node
  - Clear Node Edges: Removes all edges connected to this node
- **Right-click on empty space**:
  - Add Node Here: Creates a new node at the clicked position
  - Create Random Graph: Generates a new random graph layout
- **File operations**:
  - Save Graph...: Opens the save dialog
  - Load Graph...: Opens the load dialog
- **General options**:
  - Clear All Edges: Removes all edges while keeping nodes intact

### Mouse Controls

- **Left-click on a node**: Select it as the starting node (when not in edit mode)
- **Drag a node**: Reposition it (when in edit mode)
- **Left-click on the graph area**: Add a new node at that position (when in edit mode)
- **Right-click anywhere**: Open the context menu with node and graph operations

### Speed Control

- Use the slider to adjust the automatic execution speed

### Keyboard Shortcuts (still supported)

- **B**: Start Breadth-First Search
- **D**: Start Depth-First Search
- **Space**: Step through the algorithm
- **A**: Toggle automatic stepping
- **R**: Reset the simulation

## Features

- Interactive graph visualization
- Step-by-step algorithm execution
- Visual animation highlighting current nodes and edges
- Speed control for automatic execution
- Visualizes both BFS and DFS traversals
- Shows the current state of the Queue (BFS) or Stack (DFS)
- Displays the order of visited nodes
- User-controllable graph editing:
  - Add/remove nodes
  - Add/remove edges
  - Reposition nodes
  - Grid and snap-to-grid functionality
  - Right-click context menu for quick operations
- File operations:
  - Save graphs to JSON files
  - Load graphs from JSON files
- Fully button-based UI (no keyboard required)
- Supports both manual and automatic stepping
- Interactive feedback and help messages

## Algorithm Comparison

### Breadth-First Search (BFS)

- Uses a queue (FIFO) data structure
- Explores all neighbors at the current depth before moving to nodes at the next depth
- Finds the shortest path in an unweighted graph
- Good for finding the shortest path or minimum steps

### Depth-First Search (DFS)

- Uses a stack (LIFO) data structure
- Explores as far as possible along a branch before backtracking
- May not find the shortest path
- Good for maze solving, topological sorting, and detecting cycles

### AVL Tree

- Self-balancing binary search tree
- Maintains height balance through rotations
- Guarantees O(log n) time complexity for insert, delete, and search operations
- Height difference between left and right subtrees is at most 1
- Automatically rebalances after insertions and deletions
- Displays node values and heights for educational purposes

## AVL Tree Features

- **Interactive Operations**: Insert, delete, and search for values with visual feedback
- **Real-time Visualization**: See the tree structure update immediately after operations
- **Height Display**: Each node shows its value and height
- **Balance Visualization**: Tree automatically maintains AVL balance property
- **Search Highlighting**: Found nodes are highlighted during search operations
- **Zoom and Pan Support**: Navigate large trees easily
- **Input Validation**: Prevents invalid operations and provides user feedback

## Development

To extend this application, you might consider:

1. Adding more graph algorithms (Dijkstra's, A\*, etc.)
2. Adding support for weighted edges
3. Implementing directed graphs
4. Adding performance metrics and statistics
5. Implementing other graph theory concepts (MST, flow networks, etc.)
6. Creating visualization options for different graph layouts
7. Adding an undo/redo feature for graph edits
8. Implementing zooming and panning for larger graphs
9. Supporting export to image or PDF formats
10. Implementing an animation speed control for the traversal visualization
