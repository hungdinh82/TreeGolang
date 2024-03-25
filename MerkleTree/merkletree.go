package main

import (
	"bufio"
	"bytes" // Import gói bytes
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

var defaultTransactions = [][]byte{
	[]byte("Transaction1"),
	[]byte("Transaction2"),
}

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

	// Thêm các giao dịch mặc định vào cây
	nodes = append(nodes, createNodes(defaultTransactions)...)

	// Tạo các nút mới nếu có yêu cầu thêm nút từ người dùng
	if len(data) > 0 {
		nodes = append(nodes, createNodes(data)...)
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

// createNodes tạo các nút từ slice dữ liệu
func createNodes(data [][]byte) []*Node {
	var nodes []*Node
	for _, d := range data {
		nodes = append(nodes, NewNode(d))
	}
	return nodes
}

// FindTransactionData tìm kiếm và trả về dữ liệu của một giao dịch cụ thể trong cây Merkle
func FindTransactionData(root *Node, data []byte) []byte {
	// Nếu cây rỗng hoặc không có giao dịch nào, trả về nil
	if root == nil || len(data) == 0 {
		return nil
	}

	// Nếu nút hiện tại chứa dữ liệu của giao dịch, trả về dữ liệu đó
	if bytes.Equal(root.Data, data) {
		return root.Data
	}

	// Tìm kiếm trong cây con bên trái
	leftData := FindTransactionData(root.Left, data)
	if leftData != nil {
		return leftData
	}

	// Tìm kiếm trong cây con bên phải
	rightData := FindTransactionData(root.Right, data)
	if rightData != nil {
		return rightData
	}

	// Nếu không tìm thấy dữ liệu trong cây, trả về nil
	return nil
}

func main() {
	data := [][]byte{
		[]byte("Transaction1"),
		[]byte("Transaction2"),
	}

	// Tạo cây Merkle từ dữ liệu
	merkleTree := NewMerkleTree(data)

	// In ra cây Merkle
	fmt.Println("Merkle Tree:")
	PrintMerkleTree(merkleTree.Root, 0)

	// Nhập số lượng nút mới từ người dùng
	fmt.Print("Nhập số lượng nút mới cần thêm vào cây: ")
	var numNodes int
	fmt.Scanln(&numNodes)

	// Khởi tạo slice để lưu trữ dữ liệu các nút mới
	newData := make([][]byte, numNodes)

	// Nhập dữ liệu cho các nút mới từ người dùng
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i < numNodes; i++ {
		fmt.Printf("Nhập dữ liệu cho nút %d: ", i+1)
		scanner.Scan()
		data := scanner.Bytes()
		newData[i] = data
	}

	// Tạo cây Merkle mới từ dữ liệu hiện có và dữ liệu mới
	merkleTree = NewMerkleTree(newData)

	// In ra cây Merkle sau khi thêm nút mới
	fmt.Println("Merkle Tree sau khi thêm nút mới:")
	PrintMerkleTree(merkleTree.Root, 0)

	// Nhập dữ liệu giao dịch từ người dùng
	var transactionData string
	fmt.Print("Nhập dữ liệu của giao dịch để tìm: ")
	fmt.Scanln(&transactionData)

	// Chuyển dữ liệu của giao dịch từ kiểu string sang kiểu byte slice
	transactionDataBytes := []byte(transactionData)

	// Tìm nút chứa dữ liệu của giao dịch
	foundData := FindTransactionData(merkleTree.Root, transactionDataBytes)

	// Kiểm tra xem liệu giao dịch đã được tìm thấy hay không
	if foundData != nil {
		fmt.Printf("Giao dịch %s được tìm thấy trong cây Merkle.\n", transactionData)
		// In ra dữ liệu của giao dịch
		fmt.Printf("Dữ liệu của giao dịch: %s\n", string(foundData))
	} else {
		fmt.Printf("Giao dịch %s không tồn tại trong cây Merkle.\n", transactionData)
	}
}
