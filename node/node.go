// package node defines a standard node interface for information stored in the zet
// each node type provides methods to access the node's data
package node

import (
	"os"
	"strings"
)

type Node interface {
	Title() string
	Load(path string) error
	path() string
}

type MarkdownNode struct {
	path string
}

func (n *MarkdownNode) Load(path string) error {
	// check the file exists
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	n.path = path
	return nil
}

func (n *MarkdownNode) read() ([]byte, error) {
	// read the file and return the contents
	contents, err := os.ReadFile(n.path)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (n MarkdownNode) Title() string {
	// Read the file from the path
	c, err := n.read()
	if err != nil && len(c) == 0 {
		return ""
	}

	contents := string(c)
	// Find the first line that starts with a hash
	// and return the text after the hash
	for _, line := range strings.Split(string(contents), "\n") {
		if strings.HasPrefix(line, "#") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	return ""

}
