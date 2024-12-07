package algs

import (
	"github.com/realrabbithouse/go-play/comparable"
)

type TreeNode struct {
	left  *TreeNode
	right *TreeNode
	key   comparable.Comparable
	value any
	n     int // number of nodes in subtree
}

func (n TreeNode) KV() (comparable.Comparable, any) {
	return n.key, n.value
}

func (n TreeNode) Left() *TreeNode {
	return n.left
}

func (n TreeNode) Right() *TreeNode {
	return n.right
}

func NewTreeNode(key comparable.Comparable, value any) *TreeNode {
	return &TreeNode{key: key, value: value, n: 1}
}

type BST struct {
	root *TreeNode
}

func (t *BST) Size() int {
	return size(t.root)
}

func size(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return node.n
}

func (t *BST) Contains(key comparable.Comparable) bool {
	return contains(t.root, key)
}

func contains(node *TreeNode, key comparable.Comparable) bool {
	if node == nil {
		return false
	}
	cmp := key.CompareTo(node.key)
	if cmp < 0 {
		return contains(node.left, key)
	} else if cmp > 0 {
		return contains(node.right, key)
	} else {
		return true
	}
}

func (t *BST) Min() *TreeNode {
	if t.root == nil {
		return nil
	}
	cur := t.root
	for cur.left != nil {
		cur = cur.left
	}
	return cur
}

func (t *BST) DeleteMin() {
	if t.root == nil {
		return
	}
	t.root = deleteMin(t.root)
}

// deleteMin removes the minimum node from the subtree rooted at the given node.
// It returns the new root of the subtree after the minimum node has been removed.
//
// Parameters:
//   - node: A pointer to the root of the subtree from which the minimum node is to be removed.
//
// Returns:
//   - A pointer to the new root of the subtree after the minimum node has been removed.
func deleteMin(node *TreeNode) *TreeNode {
	if node.left == nil {
		return node.right
	}
	node.left = deleteMin(node.left)
	node.n = size(node.left) + size(node.right) + 1
	return node
}

func (t *BST) Max() *TreeNode {
	if t.root == nil {
		return nil
	}
	cur := t.root
	for cur.right != nil {
		cur = cur.right
	}
	return cur
}

func (t *BST) DeleteMax() {
	if t.root == nil {
		return
	}
	t.root = deleteMax(t.root)
}

// deleteMax removes the maximum node from the subtree rooted at the given node.
// It returns the new root of the subtree after the maximum node has been removed.
//
// Parameters:
//   - node: A pointer to the root of the subtree from which the maximum node is to be removed.
//
// Returns:
//   - A pointer to the new root of the subtree after the maximum node has been removed.
func deleteMax(node *TreeNode) *TreeNode {
	if node.right == nil {
		return node.left
	}
	node.right = deleteMax(node.right)
	node.n = size(node.left) + size(node.right) + 1
	return node
}
func (t *BST) Get(key comparable.Comparable) any {
	return get(t.root, key)
}

func get(node *TreeNode, key comparable.Comparable) any {
	if node == nil {
		return nil
	}
	cmp := key.CompareTo(node.key)
	if cmp < 0 {
		return get(node.left, key)
	} else if cmp > 0 {
		return get(node.right, key)
	} else {
		return node.value
	}
}

func (t *BST) Put(key comparable.Comparable, value any) {
	t.root = put(t.root, key, value)
}

func put(node *TreeNode, key comparable.Comparable, value any) *TreeNode {
	if node == nil {
		return NewTreeNode(key, value)
	}
	cmp := key.CompareTo(node.key)
	if cmp < 0 {
		node.left = put(node.left, key, value)
	} else if cmp > 0 {
		node.right = put(node.right, key, value)
	} else {
		node.value = value
	}
	node.n = size(node.left) + size(node.right) + 1
	return node
}

func (t *BST) Delete(key comparable.Comparable) {
	t.root = deleteKey(t.root, key)
}

// deleteKey removes the node with the specified key from the subtree rooted at the given node.
// It returns the new root of the subtree after the node with the specified key has been removed.
//
// Parameters:
//   - node: A pointer to the root of the subtree from which the node with the specified key is to be removed.
//   - key: The key of the node to be removed.
//
// Returns:
//   - A pointer to the new root of the subtree after the node with the specified key has been removed.
func deleteKey(node *TreeNode, key comparable.Comparable) *TreeNode {
	if node == nil {
		return nil
	}

	if key.CompareTo(node.key) < 0 {
		// key is less than node.key, so search in the left subtree
		node.left = deleteKey(node.left, key)
	} else if key.CompareTo(node.key) > 0 {
		// key is greater than node.key, so search in the right subtree
		node.right = deleteKey(node.right, key)
	} else {
		if node.left == nil {
			return node.right
		}
		if node.right == nil {
			return node.left
		}
		// find the minimum node in the right subtree and replace the current node with it
		t := node
		node = minTreeNode(t.right)
		node.right = deleteMin(t.right)
	}

	node.n = size(node.left) + size(node.right) + 1
	return node
}

func minTreeNode(node *TreeNode) *TreeNode {
	if node.left == nil {
		return node
	}
	return minTreeNode(node.left)
}

func (t *BST) Choose(i int) *TreeNode {
	n := t.Size()
	if i < 0 || i >= n {
		return nil
	}
	return choose(t.root, i)
}

// choose returns the i-th smallest node in the subtree rooted at the given node.
//
// Parameters:
//   - node: A pointer to the root of the subtree in which to find the i-th smallest node.
//   - i: The index (0-based) of the node to find.
//
// Returns:
//   - A pointer to the i-th smallest node in the subtree, or nil if the index is out of bounds
//     or the subtree is empty.
func choose(node *TreeNode, i int) *TreeNode {
	if node == nil {
		return nil
	}
	sz := size(node.left)
	if i < sz {
		return choose(node.left, i)
	} else if i > sz {
		return choose(node.right, i-sz-1)
	} else {
		return node
	}
}

func (t *BST) Rank(key comparable.Comparable) int {
	return rank(t.root, key)
}

// rank returns the number of keys in the subtree rooted at the given node that are less than the specified key.
//
// Parameters:
//   - node: A pointer to the root of the subtree in which to count the number of keys less than the specified key.
//   - key: The key to compare against the keys in the subtree.
//
// Returns:
//   - An integer representing the number of keys in the subtree that are less than the specified key.
func rank(node *TreeNode, key comparable.Comparable) int {
	if node == nil {
		return 0
	}
	cmp := key.CompareTo(node.key)
	if cmp < 0 {
		return rank(node.left, key)
	} else if cmp > 0 {
		return rank(node.right, key) + size(node.left) + 1
	} else {
		return size(node.left)
	}
}
