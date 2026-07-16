// Package build holds build-time version information injected by GoReleaser
// via -ldflags. These values are also used by the version command and
// error-reporting.
package build

// Version is the semantic version of the release (e.g. "v0.1.0").
// Overridden by GoReleaser at build time.
var Version = "dev"

// Commit is the git commit SHA from which the binary was built.
var Commit = "none"

// Date is the UTC build timestamp in RFC3339 format.
var Date = "unknown"
