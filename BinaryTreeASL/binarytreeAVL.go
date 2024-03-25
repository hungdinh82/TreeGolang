package main

import (
	"fmt"
	"math"
)

type Node struct {
	Value  int
	Left   *Node
	Right  *Node
	Height int
}

func newNode(value int) *Node {
	return &Node{Value: value, Height: 1}
}

func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func balanceFactor(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func rotateRight(y *Node) *Node {
	x := y.Left
	t := x.Right

	x.Right = y
	y.Left = t

	y.Height = 1 + int(math.Max(float64(height(y.Left)), float64(height(y.Right))))
	x.Height = 1 + int(math.Max(float64(height(x.Left)), float64(height(x.Right))))

	return x
}

func rotateLeft(x *Node) *Node {
	y := x.Right
	t := y.Left

	y.Left = x
	x.Right = t

	x.Height = 1 + int(math.Max(float64(height(x.Left)), float64(height(x.Right))))
	y.Height = 1 + int(math.Max(float64(height(y.Left)), float64(height(y.Right))))

	return y
}

func insert(root *Node, value int) *Node {
	if root == nil {
		return newNode(value)
	}

	if value < root.Value {
		root.Left = insert(root.Left, value)
	} else if value > root.Value {
		root.Right = insert(root.Right, value)
	}

	root.Height = 1 + int(math.Max(float64(height(root.Left)), float64(height(root.Right))))

	balance := balanceFactor(root)

	// Trường hợp mất cân bằng trái-trái
	if balance > 1 && value < root.Left.Value {
		return rotateRight(root)
	}
	// Trường hợp mất cân bằng phải-phải
	if balance < -1 && value > root.Right.Value {
		return rotateLeft(root)
	}
	// Trường hợp mất cân bằng trái-phải  thì sẽ bẻ trái để thành => trái trái rồi bẻ phải để cân bằng
	if balance > 1 && value > root.Left.Value {
		root.Left = rotateLeft(root.Left)
		return rotateRight(root)
	}
	// Trường hợp mất cân bằng phải-trái  thì sẽ bẻ phải để thành => phải phải rồi bẻ trái để cân bằng
	if balance < -1 && value < root.Right.Value {
		root.Right = rotateRight(root.Right)
		return rotateLeft(root)
	}

	return root
}

func inorderTraversal(root *Node) {
	if root != nil {
		inorderTraversal(root.Left)
		fmt.Printf("Value: %d, Height: %d\n", root.Value, root.Height)
		inorderTraversal(root.Right)
	}
}

func main() {
	var root *Node
	values := []int{30, 20, 10, 15, 40, 25, 27, 26, 5, 13, 14}

	for _, value := range values {
		root = insert(root, value)
	}

	fmt.Println("Inorder traversal of the constructed AVL tree:")
	inorderTraversal(root)
	fmt.Println()
}
