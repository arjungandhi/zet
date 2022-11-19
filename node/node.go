// Package node is a package that contains the Node interface and its implementations
package node

type Node interface {
	Title() string
	Load(path string) error
	path() string
}
