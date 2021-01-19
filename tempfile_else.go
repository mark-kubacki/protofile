// +build !linux

package protofile

import (
	"io/ioutil"
	"os"
)

// CreateTemp wraps ioutil.TempFile(dir, ".*")
// on this operating system.
//
// File name and dir won't be the same.
func CreateTemp(dir string) (*os.File, error) {
	return ioutil.TempFile(dir, ".*")
}

// IsTempfileNotSupported always returns 'false'
// on this operating system.
func IsTempfileNotSupported(err error) bool {
	return false
}
