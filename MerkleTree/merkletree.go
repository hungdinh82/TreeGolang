package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Node đại diện cho một nút trong cây Merkle
type Node struct {
	Data  []byte
	Left  *Node
	Right *Node
}

// NewNode tạo một nút mới với dữ liệu nhất định
func NewNode(data []byte) *Node {
	return &Node{Data: data}
}

// MerkleTree đại diện cho một cây Merkle với nút gốc và chiều cao
type MerkleTree struct {
	Root  *Node
	Depth int
}

// NewMerkleTree tạo một cây Merkle từ danh sách các phần tử
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*Node

	// Tạo các nút lá từ dữ liệu
	for _, d := range data {
		nodes = append(nodes, NewNode(d))
	}

	for len(nodes) > 1 {
		var level []*Node

		// Kiểm tra xem số lượng nút là số lẻ hay chẵn
		if len(nodes)%2 != 0 {
			// Thêm một nút lá giả với dữ liệu là nút cuối cùng
			lastNode := NewNode(nodes[len(nodes)-1].Data)
			nodes = append(nodes, lastNode)
		}

		for i := 0; i < len(nodes); i += 2 {
			node1 := nodes[i]
			node2 := nodes[i+1]

			// Tạo một nút mới với hai nút con
			newNode := NewNode(append(node1.Data, node2.Data...))
			newNode.Left = node1
			newNode.Right = node2

			// Tính toán mã băm của nút mới
			hash := sha256.Sum256(newNode.Data)
			newNode.Data = hash[:]

			level = append(level, newNode)
		}

		nodes = level
	}

	return &MerkleTree{Root: nodes[0], Depth: calcDepth(len(data))}
}

// calcDepth tính toán chiều cao của cây Merkle dựa trên số lượng phần tử
func calcDepth(numElements int) int {
	if numElements == 0 {
		return 0
	}
	depth := 1
	for numElements > 1 {
		numElements = (numElements + 1) / 2
		depth++
	}
	return depth
}

// PrintMerkleTree in ra cây Merkle
func PrintMerkleTree(node *Node, depth int) {
	if node == nil {
		return
	}
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(hex.EncodeToString(node.Data))
	PrintMerkleTree(node.Left, depth+1)
	PrintMerkleTree(node.Right, depth+1)
}

func main() {
	data := [][]byte{
		[]byte("Transaction1"),
		[]byte("Transaction2"),
		[]byte("Transaction3"),
		[]byte("Transaction4"),
		
	}

	// Tạo cây Merkle từ dữ liệu
	merkleTree := NewMerkleTree(data)

	// In ra cây Merkle
	fmt.Println("Merkle Tree:")
	PrintMerkleTree(merkleTree.Root, 0)
}
