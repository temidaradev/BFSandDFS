package algorithms

// TraversalMode defines whether we're using BFS, DFS, or idle
type TraversalMode int

const (
	ModeIdle TraversalMode = iota
	ModeBFS
	ModeDFS
)

// BFSStep performs one step of the BFS algorithm
// It takes the current queue, visited map, and node list
// Returns the updated queue, a newly visited node (if any), and whether the algorithm is done
func BFSStep(queue []int, visited map[int]bool, neighbors map[int][]int) ([]int, int, bool) {
	if len(queue) == 0 {
		return queue, -1, true // Algorithm is done
	}

	// Dequeue the first node (FIFO)
	n := queue[0]
	queue = queue[1:]

	// If already visited, continue to the next step
	if visited[n] {
		return queue, -1, false
	}

	// Mark as visited and add unvisited neighbors to the queue
	for _, nb := range neighbors[n] {
		if !visited[nb] {
			queue = append(queue, nb)
		}
	}

	return queue, n, false
}

// DFSStep performs one step of the DFS algorithm
// It takes the current stack, visited map, and node list
// Returns the updated stack, a newly visited node (if any), and whether the algorithm is done
func DFSStep(stack []int, visited map[int]bool, neighbors map[int][]int) ([]int, int, bool) {
	if len(stack) == 0 {
		return stack, -1, true // Algorithm is done
	}

	// Pop the top node (LIFO)
	n := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	// If already visited, continue to the next step
	if visited[n] {
		return stack, -1, false
	}

	// Mark as visited and add unvisited neighbors to the stack
	for _, nb := range neighbors[n] {
		if !visited[nb] {
			stack = append(stack, nb)
		}
	}

	return stack, n, false
}
