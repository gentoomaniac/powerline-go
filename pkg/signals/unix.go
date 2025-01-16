//go:build !windows

package signals

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func ResolveSignalCode(signal int) string {
	return unix.SignalName(syscall.Signal(signal))
}
