package search

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"

	"github.com/arjungandhi/go-utils/pkg/shell"

	node "github.com/arjungandhi/zet/pkg/node"
)

func init() {
	Z.Vars.SoftInit()
}

var maxChars = 88
var maxLines = 2

var Cmd = &Z.Cmd{
	Name:    "search",
	Summary: "search the zettelcasten for some text",
	Usage:   "search [Query]",
	MinArgs: 0,
	MaxArgs: 1,
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(x *Z.Cmd, args ...string) error {
		initQuery := ""
		if len(args) > 0 {
			initQuery = args[0]
		}
		n, err := Search(initQuery)
		if err != nil {
			return err
		}

		fmt.Println(n.Path())

		return nil
	},
}

// Search the zettelcasten for some text
func Search(initial_query string) (node.Node, error) {
	// check for commands we need on the system
	// fzf
	fzf := shell.CheckCommand("fzf")
	if !fzf {
		return nil, fmt.Errorf("fzf not found please install it")
	}

	zetDir := Z.Vars.Get(".zet.zetdir")
	nodes := node.LoadAll(zetDir)

	search_path := []string{}
	for _, n := range nodes {
		search_path = append(search_path, n.Title())
	}

	search_str := strings.Join(search_path, "\n")

	fzf_args := []string{
		fmt.Sprintf("--query=%s", initial_query),
		`--layout=reverse`,
		`-1`,
	}

	cmd := exec.Command("fzf", fzf_args...)
	cmd.Stdin = strings.NewReader(search_str)
	cmd.Stderr = os.Stderr

	selected, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	selected_str := strings.TrimSpace(string(selected))
	title := strings.Split(selected_str, ":")[0]

	for _, n := range nodes {
		if n.Title() == title {
			return n, nil
		}
	}

	return nil, nil
}
