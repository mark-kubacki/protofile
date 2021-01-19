// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package protofile

import (
	"os"
)

// Hardlink wraps os.Link on this operating system.
func Hardlink(f *os.File, newname string) error {
	return os.Link(f.Name(), newname)
}
