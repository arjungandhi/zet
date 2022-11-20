package node_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/arjungandhi/zet/pkg/node"
)

type titleTest struct {
	path string
	want string
}

// Test Title()
var markdownTitleTests = []titleTest{
	{
		path: "testdata/bye.md",
		want: "Bye",
	},
	{
		path: "testdata/hello.md",
		want: "Hello",
	},
	{
		path: "testdata/cat.md",
		want: "Cat",
	},
	{
		path: "testdata/list.md",
		want: "List",
	},
	{
		path: "testdata/url.md",
		want: "Url",
	},
}

func TestMarkdownTitles(t *testing.T) {
	for _, test := range markdownTitleTests {
		node := node.MarkdownNode{}
		node.Load(test.path)
		if node.Title() != test.want {
			t.Errorf("Title() = %q, want %q", node.Title(), test.want)
		}
	}
}

// Test Search()
type searchTest struct {
	paths []string
	query string
	want  []string
}

var files = []string{
	"testdata/bye.md",
	"testdata/hello.md",
	"testdata/cat.md",
	"testdata/list.md",
	"testdata/url.md",
}

var searchTests = []searchTest{
	{
		paths: files,
		query: "hello",
		want: []string{
			"# Hello",
			"Hello",
			"1. Hello",
		},
	},
	{
		paths: files,
		query: "bye",
		want: []string{
			"# Bye",
			"goodbye",
		},
	},
	{
		paths: files,
		query: "cat",
		want: []string{
			"# Cat",
			"1. Cat",
		},
	},
	{
		paths: files,
		query: "list",
		want: []string{
			"# List",
			"1. List",
		},
	},
	{
		paths: files,
		query: "url",
		want: []string{
			"# Url",
			"https://url.com/",
			"1. Url",
		},
	},
}

func TestMarkdownSearch(t *testing.T) {
	for _, test := range searchTests {
		results := []string{}
		for _, path := range test.paths {
			node := node.MarkdownNode{}
			node.Load(path)
			results = append(results, node.Search(test.query)...)
		}
		// sort the results
		sort.Strings(results)
		sort.Strings(test.want)
		if !reflect.DeepEqual(results, test.want) {
			t.Errorf("Search on Files returned %q, expected %q", results, test.want)
		}
	}
}
