// Command kairos is a Git Context Engine that generates commit messages,
// tags, releases, and changelog entries using local or remote LLMs.
package main

import "github.com/afeldman/kairos/internal/cmd"

func main() {
	cmd.Execute()
}
