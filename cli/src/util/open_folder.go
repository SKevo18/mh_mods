//go:build !windows && !darwin

package util

import (
	"fmt"
)

// Opens the given folder
func OpenFolder(path string) error {
	// not implemented
	fmt.Println("opening folder is not implemented on Linux")
	return nil
}
