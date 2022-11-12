// Package Node defines the Node interface, the Node interface is the base implementation
// for any piece of knowledge in the graph
package node

import (
	"fmt"
	"strings"
)

type Node interface {
	Search(query string) bool
	Edit()
	View()
}

type BasicNode struct {
	Name string
}

func (n *BasicNode) Search(query string) bool {
	return strings.Contains(n.Name, query)
}

func (n *BasicNode) Edit() {
}

func (n *BasicNode) View() {
	fmt.Println(n.Name)
}
