package node_test

import (
	"testing"

	"github.com/arjungandhi/zet/node"
)

type test struct {
	path string
	want string
}

var markdownTitleTests = []test{
	{
		path: "testdata/test.md",
		want: "I am a Title",
	},
}

func TestMarkdownNode(t *testing.T) {
	for _, test := range markdownTitleTests {
		node := node.MarkdownNode{}
		node.Load(test.path)
		if node.Title() != test.want {
			t.Errorf("Title() = %q, want %q", node.Title(), test.want)
		}
	}
}
