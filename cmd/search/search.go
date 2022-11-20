package search

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"

	"github.com/arjungandhi/zet/pkg/node"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{
	Name:    "search",
	Summary: "search the zettelcasten for some text",
	Usage:   "search [Query]",
	MinArgs: 1,
	MaxArgs: 1,
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(x *Z.Cmd, args ...string) error {
		query := args[0]

		zetdir := Z.Vars.Get(".zet.zetdir")

		matchCount := 1
		matches := map[int]string{}
		filepath.Walk(zetdir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			switch ext {
			case ".md":
				n := node.MarkdownNode{}
				err := n.Load(path)
				if err == nil {
					search := n.Search(query)
					if len(search) > 0 {
						fmt.Printf("%d. %s\n", matchCount, n.Title())
						for _, s := range search {
							fmt.Printf("\t%s\n", s)
						}
						matches[matchCount] = path
						matchCount++
					}
				}
			}
			return nil
		})
		// after we have all the matches
		// save them to the vars
		b, err := json.Marshal(matches)
		if err != nil {
			return err
		}
		Z.Vars.Set(".zet.list", string(b))

		return nil
	},
}

func mdTitle(path string) string {
	return ""
}
