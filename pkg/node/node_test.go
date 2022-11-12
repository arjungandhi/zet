package node_test

import (
	"testing"

	"github.com/arjungandhi/zet/pkg/node"
)

func TestSearch(t *testing.T) {
	n := node.BasicNode{Name: "test"}
	if !n.Search("test") {
		t.Error("Search failed")
	}
}

func TestView(t *testing.T) {
	n := node.BasicNode{Name: "test"}
	n.View()
}

func TestEdit(t *testing.T) {
	n := node.BasicNode{Name: "test"}
	n.Edit()
}
