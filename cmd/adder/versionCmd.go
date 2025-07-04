package main

import (
	"fmt"
	"runtime/debug"

	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/spf13/cobra"
)

// getVersionInfo returns version information, preferring build-time values
// but falling back to module info when built with go install
func getVersionInfo() (string, string, string) {
	// If version was set at build time (e.g., by GoReleaser), use it
	if version != "dev" && version != "unknown" {
		return version, commit, date
	}

	// Otherwise try to get version from build info (go install)
	if info, ok := debug.ReadBuildInfo(); ok {
		moduleVersion := info.Main.Version
		var vcsRevision, vcsTime string
		
		// Extract VCS info from build settings
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				vcsRevision = setting.Value
			case "vcs.time":
				vcsTime = setting.Value
			}
		}
		
		// Use module version if available, otherwise fall back to VCS info
		if moduleVersion != "" && moduleVersion != "(devel)" {
			return moduleVersion, vcsRevision, vcsTime
		}
		
		// For development builds, show VCS info if available
		if vcsRevision != "" {
			return "dev-" + vcsRevision[:7], vcsRevision, vcsTime
		}
	}

	// Final fallback to original values
	return version, commit, date
}

// versionCmd processes the version command request to display version information.
func versionCmd(_ *cobra.Command, _ *generated.VersionRequest) error {
	ver, com, dt := getVersionInfo()
	fmt.Printf("adder version %s\n", ver)
	fmt.Printf("commit: %s\n", com)
	fmt.Printf("built at: %s\n", dt)
	return nil
}
