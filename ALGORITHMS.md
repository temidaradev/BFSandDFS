# Graph Algorithm Simulator - Extended Features

This document describes the new graph algorithms that have been added to the BFS/DFS simulator.

## Newly Added Algorithms

### 1. Dijkstra's Algorithm

**Purpose**: Finds the shortest path from a source node to all other nodes in a weighted graph.

**Visualization**:

- Shows distance labels next to each node
- Distances are displayed in red text next to nodes
- Source node is highlighted in orange

**Usage**: Click the "Dijkstra" button to run the algorithm from node 0.

**Time Complexity**: O((V + E) log V) where V is vertices and E is edges.

### 2. A\* (A-Star) Search Algorithm

**Purpose**: Finds the shortest path between two specific nodes using a heuristic function.

**Visualization**:

- The optimal path is highlighted in red
- Uses Euclidean distance as heuristic
- Finds path from first node (0) to last node

**Usage**: Click the "A\*" button to find path from node 0 to the last node.

**Time Complexity**: O(b^d) where b is branching factor and d is depth of solution.

### 3. Topological Sort

**Purpose**: Orders vertices in a directed acyclic graph (DAG) such that for every directed edge (u,v), vertex u comes before v.

**Visualization**:

- Shows order numbers next to nodes in blue
- Numbers indicate the topological ordering (1, 2, 3, etc.)

**Usage**: Click the "Topo Sort" button to compute topological ordering.

**Time Complexity**: O(V + E)

### 4. Kruskal's Algorithm (Minimum Spanning Tree)

**Purpose**: Finds a minimum spanning tree that connects all vertices with minimum total edge weight.

**Visualization**:

- MST edges are highlighted in green
- Edge weights are displayed on the selected edges
- Uses Union-Find data structure

**Usage**: Click the "Kruskal" button to compute MST using Kruskal's algorithm.

**Time Complexity**: O(E log E)

### 5. Prim's Algorithm (Minimum Spanning Tree)

**Purpose**: Alternative algorithm to find minimum spanning tree by growing the tree from a starting vertex.

**Visualization**:

- MST edges are highlighted in green
- Edge weights are displayed on the selected edges
- Uses priority queue for efficient edge selection

**Usage**: Click the "Prim" button to compute MST using Prim's algorithm.

**Time Complexity**: O((V + E) log V)

### 6. Tarjan's Algorithm (Strongly Connected Components)

**Purpose**: Finds all strongly connected components in a directed graph using DFS and low-link values.

**Visualization**:

- Each SCC is labeled with different colors
- Shows "SCC1", "SCC2", etc. labels on nodes

**Usage**: Click the "Tarjan" button to find SCCs using Tarjan's algorithm.

**Time Complexity**: O(V + E)

### 7. Kosaraju's Algorithm (Strongly Connected Components)

**Purpose**: Alternative algorithm to find strongly connected components using two DFS passes.

**Visualization**:

- Each SCC is labeled with different colors
- Shows "SCC1", "SCC2", etc. labels on nodes

**Usage**: Click the "Kosaraju" button to find SCCs using Kosaraju's algorithm.

**Time Complexity**: O(V + E)

## Enhanced Features

### Weighted Graph Support

- The graph now supports weighted edges with random weights between 1.0 and 10.0
- Edge weights are displayed for algorithms that use them (Dijkstra, A\*, Kruskal, Prim)
- Weights are shown as floating-point numbers on edge midpoints

### Improved Visualization

- Algorithm-specific color coding for results
- Distance labels for shortest path algorithms
- Path highlighting for A\* search
- MST edge highlighting for spanning tree algorithms
- SCC grouping with color-coded labels

### New UI Layout

The interface now includes an additional row of buttons for the new algorithms:

**Row 4: Advanced Graph Algorithms**

- Dijkstra (purple)
- A\* (purple)
- Topo Sort (blue)
- Kruskal (green)
- Prim (green)
- Tarjan (orange)
- Kosaraju (orange)

### Data Structures Used

1. **Priority Queue**: Used in Dijkstra's and A\* algorithms for efficient minimum extraction
2. **Union-Find**: Used in Kruskal's algorithm for cycle detection
3. **DFS Stack**: Used in topological sort and SCC algorithms
4. **Adjacency Lists**: Enhanced to support both weighted and unweighted representations

## Testing

A comprehensive test suite (`test_algorithms.go`) demonstrates all algorithms working correctly:

```bash
go run test_algorithms.go
```

This will show sample output for all algorithms on a randomly generated graph.

## Technical Implementation

### Key Files Modified/Added:

1. **`internal/algorithms/graph_algorithms.go`** - New file containing all advanced graph algorithms
2. **`internal/algorithms/traversal.go`** - Updated with new traversal modes
3. **`internal/graph/graph.go`** - Enhanced to support weighted edges
4. **`internal/simulator/simulator.go`** - Added Start methods for new algorithms
5. **`internal/ui/ui.go`** - Updated UI with new buttons and visualization methods

### Algorithm Results Storage:

- `ShortestPaths`: Dijkstra's distance results
- `Path`: A\* pathfinding results
- `MST`: Minimum spanning tree edges
- `SCCs`: Strongly connected components
- `TopOrder`: Topological ordering

## Future Enhancements

Potential areas for expansion:

- Step-by-step visualization for complex algorithms
- Interactive node selection for source/destination
- Graph editing capabilities
- Additional algorithms (Floyd-Warshall, Bellman-Ford, etc.)
- Performance benchmarking tools
- Graph import/export functionality
