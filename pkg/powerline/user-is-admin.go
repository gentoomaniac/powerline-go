//go:build !windows
// +build !windows

package powerline

import (
	"os"
)

func userIsAdmin() bool {
	return os.Getuid() == 0
}
