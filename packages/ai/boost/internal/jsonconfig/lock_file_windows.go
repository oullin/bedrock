//go:build windows

package jsonconfig

import (
	"os"

	"golang.org/x/sys/windows"
)

const (
	lockLengthLow  = ^uint32(0)
	lockLengthHigh = ^uint32(0)
)

func lockFile(path string) (unlock func(), err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)

	if err != nil {
		return nil, err
	}

	if err := lockExclusive(f); err != nil {
		f.Close()

		return nil, err
	}

	return func() {
		_ = unlockFile(f)
		_ = f.Close()
	}, nil
}

func lockExclusive(file *os.File) error {
	overlapped := new(windows.Overlapped)

	return windows.LockFileEx(
		windows.Handle(file.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK,
		0,
		lockLengthLow,
		lockLengthHigh,
		overlapped,
	)
}

func unlockFile(file *os.File) error {
	overlapped := new(windows.Overlapped)

	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, lockLengthLow, lockLengthHigh, overlapped)
}
