// Package version exposes build metadata. The variables are overridden at
// build time with -ldflags, for example:
//
//	go build -ldflags "-X github.com/aminbenmansour/terrakube-go/internal/version.Version=1.2.3"
package version

// These values are replaced via -ldflags in the Makefile and CI.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Info is the JSON representation returned by the /version endpoint.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns the current build metadata.
func Get() Info {
	return Info{Version: Version, Commit: Commit, Date: Date}
}
