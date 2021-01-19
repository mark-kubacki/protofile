// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protofile

import (
	"os"
	"sync/atomic"
	"syscall"

	"golang.org/x/sys/unix"
)

// Hardlink registers the file under newname with the file system,
// usually as hard link. Unlike os.Link this does work on the
// file descriptor.
//
// Prefer using Hardlink for files without names.
func Hardlink(f *os.File, newname string) error {
	if f == nil {
		return syscall.EINVAL
	}
	fd := int(f.Fd())
	return hardlink(fd, newname)
}

// Relates to CAP_DAC_READ_SEARCH, which usually unprivileged
// users don't have.
// Gets downgraded atomically.
var capState int32 = unix.CAP_DAC_READ_SEARCH

func hardlink(fd int, newname string) error {
	if (atomic.LoadInt32(&capState) & unix.CAP_DAC_READ_SEARCH) == 0 {
		oldpath := "/proc/self/fd/" + uitoa(uint(fd)) // 'self' usually is a symlink.
		// Given the absolute oldpath the first argument gets ignored according to docs.
		// Go with fd though, to recover from an inaccessible oldpath.
		return unix.Linkat(fd, oldpath, unix.AT_FDCWD, newname, unix.AT_SYMLINK_FOLLOW)
	}
	err := unix.Linkat(fd, "", unix.AT_FDCWD, newname, unix.AT_EMPTY_PATH)
	if err != syscall.ENOENT {
		return err
	}
	oldCap := atomic.LoadInt32(&capState)
	clearedCap := oldCap &^ unix.CAP_DAC_READ_SEARCH
	atomic.CompareAndSwapInt32(&capState, oldCap, clearedCap)
	return hardlink(fd, newname)
}

// Here to avoid importing "fmt".
func uitoa(val uint) string {
	var buf [32]byte // big enough for int64
	i := len(buf) - 1
	for val >= 10 {
		buf[i] = byte(val%10 + '0')
		i--
		val /= 10
	}
	buf[i] = byte(val + '0')
	return string(buf[i:])
}
