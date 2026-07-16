// Package detect identifies a project's primary language/toolchain from
// well-known manifest files in its root directory.
package detect

import (
	"os"
	"path/filepath"
)

// ProjectType identifies the toolchain a project uses.
type ProjectType string

const (
	TypeGo      ProjectType = "go"
	TypeNode    ProjectType = "node"
	TypeRust    ProjectType = "rust"
	TypePython  ProjectType = "python"
	TypeUnknown ProjectType = "unknown"
)

var manifests = []struct {
	file string
	typ  ProjectType
}{
	{"go.mod", TypeGo},
	{"package.json", TypeNode},
	{"Cargo.toml", TypeRust},
	{"pyproject.toml", TypePython},
}

// Detect returns the ProjectType for the project rooted at dir, based on
// the first matching manifest file found. Returns TypeUnknown if none match.
func Detect(dir string) ProjectType {
	for _, m := range manifests {
		if _, err := os.Stat(filepath.Join(dir, m.file)); err == nil {
			return m.typ
		}
	}
	return TypeUnknown
}
