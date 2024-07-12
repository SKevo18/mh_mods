package utils

import (
	"os"
	"syscall"
)

func RestartApp() error {
	// restart the app
	defer os.Exit(0)
	return syscall.Exec(os.Args[0], os.Args, os.Environ())
}
