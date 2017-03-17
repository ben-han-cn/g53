package domaintree

import (
	"g53"
)

type RBNodeFlag uint32

const (
	NF_CALLBACK RBNodeFlag = 1
	NF_USER1               = 0x80000000
)

type RBNodeColor int

const (
	BLACK RBNodeColor = 0
	RED   RBNodeColor = 1
)

func (color RBNodeColor) String() string {
	if color == BLACK {
		return "black"
	} else {
		return "red"
	}
}

type Node struct {
	parent *Node
	left   *Node
	right  *Node
	color  RBNodeColor

	down *Node
	flag RBNodeFlag
	name *g53.Name
	data interface{}
}

var NULL_NODE *Node

func init() {
	if NULL_NODE == nil {
		NULL_NODE = &Node{
			parent: NULL_NODE,
			left:   NULL_NODE,
			right:  NULL_NODE,
			color:  BLACK,
			down:   NULL_NODE,
			name:   g53.Root,
		}
	}
}

func NewNode(name *g53.Name) *Node {
	return &Node{
		parent: NULL_NODE,
		left:   NULL_NODE,
		right:  NULL_NODE,
		color:  RED,
		down:   NULL_NODE,
		name:   name,
	}
}

func (node *Node) IsEmpty() bool {
	return node.data == nil
}

func (node *Node) GetFlag(flag RBNodeFlag) bool {
	return (node.flag & flag) != 0
}

func (node *Node) SetFlag(flag RBNodeFlag, set bool) {
	if set {
		node.flag = node.flag | flag
	} else {
		node.flag = node.flag & (^flag)
	}
}

func (node *Node) successor() *Node {
	current := node
	if node.right != NULL_NODE {
		current = node.right
		for current.left != NULL_NODE {
			current = current.left
		}
		return current
	}

	// Otherwise go up until we find the first left branch on our path to
	// root.  If found, the parent of the branch is the successor.
	// Otherwise, we return the null node
	parent := current.parent
	for parent != NULL_NODE && current == parent.right {
		current = parent
		parent = parent.parent
	}
	return parent
}

func (node *Node) SetData(data interface{}) {
	node.data = data
}

func (node *Node) Data() interface{} {
	return node.data
}

func (node *Node) Name() *g53.Name {
	return node.name
}