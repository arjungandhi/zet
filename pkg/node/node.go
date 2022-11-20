// Package node is a package that contains the Node interface and its implementations
package node

type Node interface {
	// Title returns the title of the node
	Title() string
	// Search searches the node for the provided query
	// and returns a list of results
	Search(query string) []string
	// Load loads the node into the provided path
	Load(path string) error
}
