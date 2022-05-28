package main

import (
	"fmt"
	"runtime/debug"

	"go.bbkane.com/warg/command"
)

// goreleaser will fill this in.
// Can fill in manually with `go build -ldflags "-X main.version=myversion"`.
// Or run with `go run -ldflags "-X main.version=myversion" . version`
var version string

func getVersion() string {
	if version != "" {
		// embedded by goreleaser
		return version
	}
	// If installed via `go install`, we'll be able to read runtime version info
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown version: error reading build info"
	}
	// when run with `go run`, this will return "(devel)"
	return info.Main.Version
}

func printVersion(_ command.Context) error {
	fmt.Println(getVersion())
	return nil
}
