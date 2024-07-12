//go:build windows

package utils

import (
	"os/exec"
)

// Opens the given folder in the file explorer.
func OpenFolder(path string) error {
	return exec.Command("start", path).Start()
}
