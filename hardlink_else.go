// +build !linux

package protofile

import (
	"os"
)

// Hardlink wraps os.Link on this operating system.
func Hardlink(f *os.File, newname string) error {
	return os.Link(f.Name(), newname)
}
