//go:build !windows
// +build !windows

package segments

import (
	"os"
)

func userIsAdmin() bool {
	return os.Getuid() == 0
}
