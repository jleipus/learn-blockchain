package merkel

import "crypto/sha256"

type MerkelTree struct {
	Root *MerkelNode
}

type MerkelNode struct {
	Left  *MerkelNode
	Right *MerkelNode
	data  []byte
}

func newNode(left, right *MerkelNode, data []byte) *MerkelNode {
	var hash [32]byte
	if left == nil && right == nil {
		hash = sha256.Sum256(data)
	} else {
		prevHashes := append(left.data, right.data...)
		hash = sha256.Sum256(prevHashes)
	}

	return &MerkelNode{
		Left:  left,
		Right: right,
		data:  hash[:],
	}
}

func NewTree(data [][]byte) *MerkelTree {
	if len(data) == 0 {
		return &MerkelTree{
			Root: newNode(nil, nil, nil),
		}
	}

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1]) // Duplicate last element if odd count
	}

	nodes := make([]*MerkelNode, 0, len(data))
	for _, d := range data {
		node := newNode(nil, nil, d)
		nodes = append(nodes, node)
	}

	for range len(nodes) / 2 {
		var newLevel []*MerkelNode

		for j := 0; j < len(nodes); j += 2 {
			node := newNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}

		nodes = newLevel
	}

	return &MerkelTree{
		Root: nodes[0],
	}
}

func (n *MerkelNode) GetData() []byte {
	return n.data
}
