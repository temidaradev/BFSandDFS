package algorithms

// AVLNode represents a node in an AVL tree
type AVLNode struct {
	Value    int
	Left     *AVLNode
	Right    *AVLNode
	Height   int
	Parent   *AVLNode
	Position struct {
		X, Y int
	}
}

// AVLTree represents an AVL tree
type AVLTree struct {
	Root *AVLNode
}

// NewAVLTree creates a new empty AVL tree
func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getHeight returns the height of a node
func getHeight(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

// getBalance returns the balance factor of a node
func getBalance(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return getHeight(node.Left) - getHeight(node.Right)
}

// rightRotate performs a right rotation
func rightRotate(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	// Perform rotation
	x.Right = y
	y.Left = T2

	// Update heights
	y.Height = max(getHeight(y.Left), getHeight(y.Right)) + 1
	x.Height = max(getHeight(x.Left), getHeight(x.Right)) + 1

	// Update parent pointers
	if T2 != nil {
		T2.Parent = y
	}
	x.Parent = y.Parent
	y.Parent = x

	return x
}

// leftRotate performs a left rotation
func leftRotate(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	// Perform rotation
	y.Left = x
	x.Right = T2

	// Update heights
	x.Height = max(getHeight(x.Left), getHeight(x.Right)) + 1
	y.Height = max(getHeight(y.Left), getHeight(y.Right)) + 1

	// Update parent pointers
	if T2 != nil {
		T2.Parent = x
	}
	y.Parent = x.Parent
	x.Parent = y

	return y
}

// Insert adds a new value to the AVL tree
func (t *AVLTree) Insert(value int) {
	t.Root = t.insertNode(t.Root, value)
}

// insertNode recursively inserts a value into the AVL tree
func (t *AVLTree) insertNode(node *AVLNode, value int) *AVLNode {
	// Standard BST insert
	if node == nil {
		return &AVLNode{
			Value:  value,
			Height: 1,
		}
	}

	if value < node.Value {
		node.Left = t.insertNode(node.Left, value)
		node.Left.Parent = node
	} else if value > node.Value {
		node.Right = t.insertNode(node.Right, value)
		node.Right.Parent = node
	} else {
		return node // Duplicate values not allowed
	}

	// Update height
	node.Height = 1 + max(getHeight(node.Left), getHeight(node.Right))

	// Get balance factor
	balance := getBalance(node)

	// Left Left Case
	if balance > 1 && value < node.Left.Value {
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && value > node.Right.Value {
		return leftRotate(node)
	}

	// Left Right Case
	if balance > 1 && value > node.Left.Value {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Left Case
	if balance < -1 && value < node.Right.Value {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

// Delete removes a value from the AVL tree
func (t *AVLTree) Delete(value int) {
	t.Root = t.deleteNode(t.Root, value)
}

// deleteNode recursively deletes a value from the AVL tree
func (t *AVLTree) deleteNode(node *AVLNode, value int) *AVLNode {
	// Standard BST delete
	if node == nil {
		return nil
	}

	if value < node.Value {
		node.Left = t.deleteNode(node.Left, value)
	} else if value > node.Value {
		node.Right = t.deleteNode(node.Right, value)
	} else {
		// Node to delete found

		// Node with only one child or no child
		if node.Left == nil {
			temp := node.Right
			if temp != nil {
				temp.Parent = node.Parent
			}
			return temp
		} else if node.Right == nil {
			temp := node.Left
			if temp != nil {
				temp.Parent = node.Parent
			}
			return temp
		}

		// Node with two children: Get the inorder successor (smallest in right subtree)
		temp := t.getMinValueNode(node.Right)
		node.Value = temp.Value
		node.Right = t.deleteNode(node.Right, temp.Value)
	}

	if node == nil {
		return nil
	}

	// Update height
	node.Height = 1 + max(getHeight(node.Left), getHeight(node.Right))

	// Get balance factor
	balance := getBalance(node)

	// Left Left Case
	if balance > 1 && getBalance(node.Left) >= 0 {
		return rightRotate(node)
	}

	// Left Right Case
	if balance > 1 && getBalance(node.Left) < 0 {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && getBalance(node.Right) <= 0 {
		return leftRotate(node)
	}

	// Right Left Case
	if balance < -1 && getBalance(node.Right) > 0 {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

// getMinValueNode returns the node with minimum value in the tree
func (t *AVLTree) getMinValueNode(node *AVLNode) *AVLNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// Search looks for a value in the AVL tree
func (t *AVLTree) Search(value int) *AVLNode {
	return t.searchNode(t.Root, value)
}

// searchNode recursively searches for a value in the AVL tree
func (t *AVLTree) searchNode(node *AVLNode, value int) *AVLNode {
	if node == nil || node.Value == value {
		return node
	}

	if value < node.Value {
		return t.searchNode(node.Left, value)
	}
	return t.searchNode(node.Right, value)
}

// UpdatePositions updates the visual positions of all nodes in the tree
func (t *AVLTree) UpdatePositions(startX, startY, levelHeight int) {
	t.updateNodePositions(t.Root, startX, startY, levelHeight, 0)
}

// updateNodePositions recursively updates node positions for visualization
func (t *AVLTree) updateNodePositions(node *AVLNode, x, y, levelHeight, level int) {
	if node == nil {
		return
	}

	// Calculate horizontal spacing based on level
	spacing := 1 << (level + 2) // 2^(level+2)

	// Update current node position
	node.Position.X = x
	node.Position.Y = y

	// Update left subtree positions
	if node.Left != nil {
		t.updateNodePositions(node.Left, x-spacing, y+levelHeight, levelHeight, level+1)
	}

	// Update right subtree positions
	if node.Right != nil {
		t.updateNodePositions(node.Right, x+spacing, y+levelHeight, levelHeight, level+1)
	}
}
