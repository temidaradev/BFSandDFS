package algorithms

import (
	"container/heap"
	"math"
	"sort"
)

// Edge represents a weighted edge in a graph
type Edge struct {
	From, To int
	Weight   float64
}

// PriorityQueueItem represents an item in a priority queue
type PriorityQueueItem struct {
	Node     int
	Priority float64
	Index    int
}

// PriorityQueue implements a min-heap for Dijkstra's algorithm
type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityQueueItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra's Algorithm - finds shortest path from source to all other nodes
func Dijkstra(neighbors map[int][]Edge, source int, numNodes int) (map[int]float64, map[int]int) {
	dist := make(map[int]float64)
	prev := make(map[int]int)
	visited := make(map[int]bool)

	// Initialize distances
	for i := 0; i < numNodes; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}
	dist[source] = 0

	// Priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PriorityQueueItem{Node: source, Priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*PriorityQueueItem)

		if visited[current.Node] {
			continue
		}
		visited[current.Node] = true

		for _, edge := range neighbors[current.Node] {
			if visited[edge.To] {
				continue
			}

			newDist := dist[current.Node] + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				prev[edge.To] = current.Node
				heap.Push(pq, &PriorityQueueItem{Node: edge.To, Priority: newDist})
			}
		}
	}

	return dist, prev
}

// AStar - A* search algorithm with heuristic
type Position struct {
	X, Y int
}

func AStar(neighbors map[int][]Edge, start, goal int, positions map[int]Position) ([]int, float64) {
	openSet := &PriorityQueue{}
	heap.Init(openSet)

	gScore := make(map[int]float64)
	fScore := make(map[int]float64)
	cameFrom := make(map[int]int)

	for i := range neighbors {
		gScore[i] = math.Inf(1)
		fScore[i] = math.Inf(1)
	}

	gScore[start] = 0
	fScore[start] = heuristic(positions[start], positions[goal])

	heap.Push(openSet, &PriorityQueueItem{Node: start, Priority: fScore[start]})

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*PriorityQueueItem).Node

		if current == goal {
			// Reconstruct path
			path := []int{}
			for current != -1 {
				path = append([]int{current}, path...)
				if prev, exists := cameFrom[current]; exists {
					current = prev
				} else {
					break
				}
			}
			return path, gScore[goal]
		}

		for _, edge := range neighbors[current] {
			tentativeGScore := gScore[current] + edge.Weight

			if tentativeGScore < gScore[edge.To] {
				cameFrom[edge.To] = current
				gScore[edge.To] = tentativeGScore
				fScore[edge.To] = gScore[edge.To] + heuristic(positions[edge.To], positions[goal])

				heap.Push(openSet, &PriorityQueueItem{Node: edge.To, Priority: fScore[edge.To]})
			}
		}
	}

	return nil, math.Inf(1) // No path found
}

// Heuristic function for A* (Euclidean distance)
func heuristic(a, b Position) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// TopologicalSort - sorts vertices in topological order (for DAG)
func TopologicalSort(neighbors map[int][]int, numNodes int) []int {
	visited := make(map[int]bool)
	stack := []int{}

	var dfs func(int)
	dfs = func(node int) {
		visited[node] = true
		for _, neighbor := range neighbors[node] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
		stack = append([]int{node}, stack...)
	}

	for i := 0; i < numNodes; i++ {
		if !visited[i] {
			dfs(i)
		}
	}

	return stack
}

// Kruskal's Algorithm - finds minimum spanning tree
func Kruskal(edges []Edge, numNodes int) []Edge {
	// Sort edges by weight
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	// Initialize union-find
	parent := make([]int, numNodes)
	rank := make([]int, numNodes)
	for i := 0; i < numNodes; i++ {
		parent[i] = i
		rank[i] = 0
	}

	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) bool {
		px, py := find(x), find(y)
		if px == py {
			return false
		}
		if rank[px] < rank[py] {
			px, py = py, px
		}
		parent[py] = px
		if rank[px] == rank[py] {
			rank[px]++
		}
		return true
	}

	mst := []Edge{}
	for _, edge := range edges {
		if union(edge.From, edge.To) {
			mst = append(mst, edge)
			if len(mst) == numNodes-1 {
				break
			}
		}
	}

	return mst
}

// Prim's Algorithm - finds minimum spanning tree
func Prim(neighbors map[int][]Edge, numNodes int) []Edge {
	if numNodes == 0 {
		return []Edge{}
	}

	visited := make(map[int]bool)
	mst := []Edge{}
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Start with node 0
	visited[0] = true
	for _, edge := range neighbors[0] {
		heap.Push(pq, &PriorityQueueItem{Node: edge.To, Priority: edge.Weight})
	}

	for pq.Len() > 0 && len(mst) < numNodes-1 {
		item := heap.Pop(pq).(*PriorityQueueItem)
		node := item.Node
		weight := item.Priority

		if visited[node] {
			continue
		}

		visited[node] = true

		// Find the edge that led to this node with the minimum weight
		var minEdge Edge
		found := false
		for visitedNode := range visited {
			if visitedNode == node {
				continue
			}
			for _, edge := range neighbors[visitedNode] {
				if edge.To == node && edge.Weight == weight {
					minEdge = Edge{From: visitedNode, To: node, Weight: weight}
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if found {
			mst = append(mst, minEdge)
		}

		// Add new edges to priority queue
		for _, edge := range neighbors[node] {
			if !visited[edge.To] {
				heap.Push(pq, &PriorityQueueItem{Node: edge.To, Priority: edge.Weight})
			}
		}
	}

	return mst
}

// Tarjan's Algorithm - finds strongly connected components
func Tarjan(neighbors map[int][]int, numNodes int) [][]int {
	index := 0
	stack := []int{}
	onStack := make(map[int]bool)
	indices := make(map[int]int)
	lowlinks := make(map[int]int)
	sccs := [][]int{}

	var strongConnect func(int)
	strongConnect = func(v int) {
		indices[v] = index
		lowlinks[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range neighbors[v] {
			if _, exists := indices[w]; !exists {
				strongConnect(w)
				if lowlinks[w] < lowlinks[v] {
					lowlinks[v] = lowlinks[w]
				}
			} else if onStack[w] {
				if indices[w] < lowlinks[v] {
					lowlinks[v] = indices[w]
				}
			}
		}

		if lowlinks[v] == indices[v] {
			scc := []int{}
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			sccs = append(sccs, scc)
		}
	}

	for i := 0; i < numNodes; i++ {
		if _, exists := indices[i]; !exists {
			strongConnect(i)
		}
	}

	return sccs
}

// Kosaraju's Algorithm - alternative for strongly connected components
func Kosaraju(neighbors map[int][]int, numNodes int) [][]int {
	visited := make(map[int]bool)
	stack := []int{}

	// First DFS to fill stack
	var dfs1 func(int)
	dfs1 = func(v int) {
		visited[v] = true
		for _, w := range neighbors[v] {
			if !visited[w] {
				dfs1(w)
			}
		}
		stack = append(stack, v)
	}

	for i := 0; i < numNodes; i++ {
		if !visited[i] {
			dfs1(i)
		}
	}

	// Create transpose graph
	transpose := make(map[int][]int)
	for i := 0; i < numNodes; i++ {
		transpose[i] = []int{}
	}
	for from, edges := range neighbors {
		for _, to := range edges {
			transpose[to] = append(transpose[to], from)
		}
	}

	// Second DFS on transpose
	visited = make(map[int]bool)
	sccs := [][]int{}

	var dfs2 func(int, []int) []int
	dfs2 = func(v int, component []int) []int {
		visited[v] = true
		component = append(component, v)
		for _, w := range transpose[v] {
			if !visited[w] {
				component = dfs2(w, component)
			}
		}
		return component
	}

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !visited[v] {
			scc := dfs2(v, []int{})
			sccs = append(sccs, scc)
		}
	}

	return sccs
}
