package version

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// Info returns formatted version information
func Info() string {
	return "confighandler " + Version + " (built " + BuildTime + ", commit " + GitCommit + ")"
}
