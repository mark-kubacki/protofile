// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protofile

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// CreateTemp creates a nameless file in the directory that gets
// discarded unless hardlinked.
//
// Although the apparent file name is the directory name,
// do not rely on this as it's subject to change.
//
// Fallback to os.CreateTemp on syscall.EISDIR, syscall.ENOENT, or syscall.EOPNOTSUPP.
func CreateTemp(dir string) (*os.File, error) {
	// fs.FileMode(0600) is what os.CreateTemp uses.
	return os.OpenFile(dir, os.O_WRONLY|unix.O_TMPFILE|syscall.O_CLOEXEC, 0600)
}

// IsTempfileNotSupported is true for errors CreateTemp will yield
// if O_TEMPFILE is not supported by the kernel or file system.
func IsTempfileNotSupported(err error) bool {
	switch err {
	case syscall.EISDIR, syscall.ENOENT, syscall.EOPNOTSUPP:
		return true
	}
	return false
}
