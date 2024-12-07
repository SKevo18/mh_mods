//go:build darwin
package util

import (
	"os/exec"
)

// Opens the given folder in Finder
func OpenFolder(path string) error {
	return exec.Command("open", path).Start()
}
