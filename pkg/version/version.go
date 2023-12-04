package version

var (
	// Version is the tag version from the environment
	// Will be set during build process via GoReleaser
	// See also: https://pkg.go.dev/cmd/link
	Version string
)
