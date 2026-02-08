package version

import (
	"runtime/debug"
	"strings"
)

var (
	// These variables can be overridden at build time using -ldflags
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// Info returns formatted version information.
// It prefers values injected via -ldflags, but falls back to module build info
// (useful when installing via `go install github.com/...@vX.Y.Z`).
func Info() string {
	ver := Version
	// If Version looks like the development default, try to read module info
	if ver == "dev" || strings.HasPrefix(ver, "(devel)") {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				ver = info.Main.Version
			}
		}
	}

	bt := BuildTime
	gc := GitCommit

	// Build the output string
	out := "confighandler " + ver
	out += " (built " + bt
	if gc != "unknown" && gc != "" {
		out += ", commit " + gc
	}
	out += ")"
	return out
}
