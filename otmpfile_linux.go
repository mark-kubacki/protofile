// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package protofile // import "blitznote.com/src/protofile"

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func init() {
	IntentNew = intentNewUnix
}

// unixProtoFile is the variant that utilizes O_TMPFILE.
// Although it might seem as if we write data to a directory,
// it actually goes to a nameless file in said directory.
// The file gets discarded by the OS should it remain without a name but get closed.
type unixProtoFile ProtoFile

func intentNewUnix(path, filename string) (ProtoFileBehaver, error) {
	err := os.MkdirAll(path, permBitsDir)
	if err != nil {
		return nil, err
	}
	t, err := os.OpenFile(path, os.O_WRONLY|unix.O_TMPFILE|syscall.O_CLOEXEC, permBitsFile)
	// did it fail because…
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok {
			return nil, err
		}
		switch perr.Err {
		case syscall.EISDIR, syscall.ENOENT: // … kernel does not know O_TMPFILE
			// If so, don't try it again.
			IntentNew = intentNewUnixDotted
			fallthrough
		case syscall.EOPNOTSUPP: // … O_TMPFILE is not supported on this FS
			return intentNewUnixDotted(path, filename)
		default: // … something 'regular'.
			return nil, err
		}
	}

	return &unixProtoFile{
		File:      t,
		finalName: path + "/" + filename,
	}, nil
}

func (p unixProtoFile) Zap() error {
	// NOP because O_TMPFILE files that have not been named get discarded anyway.
	return p.File.Close()
}

// Persist gives the file a name.
//
// Nameless files can be identified using tuple (PID, FD) and named
// by linking the FD to a name in the filesystem on which it had been opened.
func (p unixProtoFile) Persist() error {
	err := p.File.Sync()
	if err != nil {
		return err
	}

	fd := int(p.File.Fd())
	oldpath := "/proc/self/fd/" + uitoa(uint(fd))
	err = unix.Linkat(fd, oldpath, unix.AT_FDCWD, p.finalName, unix.AT_SYMLINK_FOLLOW)
	if os.IsExist(err) { // Someone claimed our name!
		finfo, err2 := os.Stat(p.finalName)
		if err2 == nil && !finfo.IsDir() {
			os.Remove(p.finalName) // To emulate the behaviour of Create we will "overwrite" the other file.
			err = unix.Linkat(fd, oldpath, unix.AT_FDCWD, p.finalName, unix.AT_SYMLINK_FOLLOW)
		}
	}
	// 'linkat' catches many of the errors 'os.Create' would throw,
	// only with O_TMPFILE at a later point in the file's lifecycle.
	if err != nil {
		return err
	}
	p.persisted = true
	return p.Close()
}

func (p unixProtoFile) SizeWillBe(numBytes uint64) error {
	if numBytes <= reserveFileSizeThreshold {
		return nil
	}

	fd := int(p.File.Fd())
	if numBytes <= maxInt64 {
		err := syscall.Fallocate(fd, 0, 0, int64(numBytes))
		if err == syscall.EOPNOTSUPP {
			return nil
		}

		_ = unix.Fadvise(fd, 0, int64(numBytes), unix.FADV_WILLNEED)
		_ = unix.Fadvise(fd, 0, int64(numBytes), unix.FADV_SEQUENTIAL)
		return err
	}

	// Yes, every Exbibyte counts.
	err := syscall.Fallocate(fd, 0, 0, maxInt64)
	if err == syscall.EOPNOTSUPP {
		return nil
	}
	if err != nil {
		return err
	}

	err = syscall.Fallocate(fd, 0, maxInt64, int64(numBytes-maxInt64))
	if err != nil {
		return err
	}

	// These are best-efford, so we don't care about any errors.
	// For very large files this is not optimal, but covers most of use-cases for now.
	_ = unix.Fadvise(fd, 0, maxInt64, unix.FADV_WILLNEED)
	_ = unix.Fadvise(fd, 0, maxInt64, unix.FADV_SEQUENTIAL)
	return err
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
