//go:build !windows

package filesystem

import (
	"os"

	"golang.org/x/sys/unix"
)

func lockShared(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_SH)
}

func lockExclusive(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_EX)
}

func unlockFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
