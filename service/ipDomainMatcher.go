package service

import (
	"strings"
)

type Node struct {
	children map[string]*Node
	value    bool
}

func NewNode() *Node {
	return &Node{
		children: make(map[string]*Node),
		value:    false,
	}
}

func (n *Node) Insert(domain string) {
	parts := strings.Split(domain, ".")
	node := n
	for _, part := range parts {
		if _, ok := node.children[part]; !ok {
			node.children[part] = NewNode()
		}
		node = node.children[part]
	}
	node.value = true
}

func (n *Node) Match(domain string) bool {
	parts := strings.Split(domain, ".")
	node := n
	for i, part := range parts {
		if child, ok := node.children[part]; ok {
			node = child
		} else if child, ok := node.children["*"]; ok && i == len(parts)-1 {
			node = child
		} else {
			return false
		}
		if node.value {
			return true
		}
	}
	return false
}
