//go:build windows

package filesystem

import (
	"os"

	"golang.org/x/sys/windows"
)

const (
	lockLengthLow  = ^uint32(0)
	lockLengthHigh = ^uint32(0)
)

func lockShared(file *os.File) error {
	return lockFileRegion(file, 0)
}

func lockExclusive(file *os.File) error {
	return lockFileRegion(file, windows.LOCKFILE_EXCLUSIVE_LOCK)
}

func unlockFile(file *os.File) error {
	overlapped := new(windows.Overlapped)

	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, lockLengthLow, lockLengthHigh, overlapped)
}

func lockFileRegion(file *os.File, flags uint32) error {
	overlapped := new(windows.Overlapped)

	return windows.LockFileEx(windows.Handle(file.Fd()), flags, 0, lockLengthLow, lockLengthHigh, overlapped)
}
