// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protofile // import "blitznote.com/src/protofile"

import (
	"io"
	"io/ioutil"
	"os"
)

const (
	// Following bitmasks, traditionally denoted in octal, represent the values hard-coded in Golang
	// that are used by default for creating files and directories. The 'umask' is always applied implicitly.
	permBitsDir  = os.ModePerm                      // -rwxrwxrwx
	permBitsFile = os.ModePerm & ^os.FileMode(0111) // -rw-rw-rw-
)

// ProtoFileBehaver is implemented by all variants of ProtoFile.
//
// Use this in pointers to any ProtoFile you want to utilize.
type ProtoFileBehaver interface {
	// Discards a file that has not yet been persisted/closed, else a NOP.
	Zap() error

	// Emerges the file under the initially given name into observable namespace on disk.
	// This closes the file.
	Persist() error

	io.WriteCloser
}

// ProtoFile represents a file that can be discarded or named after having been written.
// (With normal files such an committment is made ex ante, on creation.)
type ProtoFile struct {
	*os.File

	persisted bool // Has this already appeared under its final name?
	finalName string
}

// IntentNew "creates" a file which, ideally, is nameless at that point.
var IntentNew func(path, filename string) (ProtoFileBehaver, error) = intentNewUniversal

type generalizedProtoFile ProtoFile

func intentNewUniversal(path, filename string) (ProtoFileBehaver, error) {
	err := os.MkdirAll(path, permBitsDir)
	if err != nil {
		return nil, err
	}
	t, err := ioutil.TempFile(path, "."+filename)
	if err != nil {
		return nil, err
	}
	return &generalizedProtoFile{
		File:      t,
		finalName: path + "/" + filename,
	}, nil
}
