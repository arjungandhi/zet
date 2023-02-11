// Package node is a package that contains the Node interface and its implementations
package node

import (
	"os"
	"path/filepath"
)

type Node interface {
	// Title returns the title of the node
	Title() string
	// Search searches the node for the provided query
	// and returns a list of results
	Search(query string) []string
	// Load loads the node into the provided path
	Load(path string) error
	// Path returns the path of the node
	Path() string
}

func LoadAll(dir string) []Node {
	nodes := []Node{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".md":
			n := &MarkdownNode{}
			err := n.Load(path)
			if err == nil {
				nodes = append(nodes, n)
			}
		}
		return nil
	})

	return nodes
}
