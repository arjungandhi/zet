package node

import (
	"os"
	"strings"
)

type MarkdownNode struct {
	contents string
	path     string
}

func (n *MarkdownNode) Load(path string) error {
	// check the file exists
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	n.path = path
	// read the file and save them to contents
	contents, err := n.read()
	if err != nil {
		return err
	}
	n.contents = string(contents)
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

func (n *MarkdownNode) Title() string {
	// Find the first line that starts with a hash
	// and return the text after the hash
	for _, line := range strings.Split(n.contents, "\n") {
		if strings.HasPrefix(line, "#") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	return ""
}

func (n *MarkdownNode) Search(query string) []string {
	// Find all lines that contain the query
	// and return them case insensitive
	var results []string
	for _, line := range strings.Split(n.contents, "\n") {
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			results = append(results, line)
		}
	}

	return results
}

func (n *MarkdownNode) Path() string {
	return n.path
}
