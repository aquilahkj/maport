package version

import (
	"fmt"
	"runtime"
)

const (
	// Proto Proto version
	Proto = "0"
	// Major Major version
	Major = "1"
	// Minor Minor version
	Minor = "1"
)

// Version gets the version
func Version() string {
	return fmt.Sprintf("%s-%s.%s", Proto, Major, Minor)
}

// FullVersion gets the full version
func FullVersion(app string) string {
	return fmt.Sprintf("%s v%s (built w/%s)", app, Version(), runtime.Version())
}
