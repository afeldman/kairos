package formatter

import "fmt"

// errUnknownStyle creates a descriptive error for unsupported styles.
func errUnknownStyle(style string) error {
	known := make([]string, 0, len(templates))
	for k := range templates {
		known = append(known, k)
	}
	return fmt.Errorf("unknown style %q (known: %v)", style, known)
}
