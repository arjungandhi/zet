package search

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"

	"github.com/arjungandhi/zet/pkg/node"
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
						for i, s := range search {
							// if the result starts with "# " then dont print it
							// its the title
							if strings.HasPrefix(s, "# ") {
								continue
							}
							// if the results has more than 50 characters then
							// get the location of the query in the string
							// and then truncate the string to 50 characters
							// fix this when min and max are implemented
							if len(s) > maxChars {
								index := strings.Index(s, query)

								// get the nearest maxChars characters to the query
								// start at index

								start := index
								end := index + len(query)
								maxEnd := len(s) - 1

								// now we interate backwards and forwards from the query
								// until we have maxChars characters
								for end-start < maxChars {
									if start > 0 {
										start--
									}
									if end < maxEnd {
										end++
									}
								}

								s = s[start:end]

								if start > 0 {
									s = s[3:]
									s = "..." + s
								}

								if end < maxEnd {
									s = s[:len(s)-4]
									s = s + "..."
								}

							}

							fmt.Printf("\t%s\n", s)

							if i >= maxLines {
								if i < len(search)-1 {
									fmt.Println("\t...")
								}
								break
							}
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
