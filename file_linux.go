// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protofile // import "blitznote.com/src/protofile"

import (
	"os"
)

// Call this to discard the file.
// If it has already been persisted (and thereby is a 'regular' one) this will be a NOP.
func (p generalizedProtoFile) Zap() error {
	if p.persisted {
		return nil
	}
	os.RemoveAll(p.File.Name())
	return p.File.Close()
}

// Promotes a proto file to a 'regular' one, which will appear under its final name.
func (p generalizedProtoFile) Persist() error {
	defer p.File.Close() // yes, this gets called up to two times
	err := p.File.Sync()
	if err != nil {
		return err
	}
	err = os.Rename(p.File.Name(), p.finalName)
	if err != nil {
		return err
	}
	p.persisted = true
	return p.File.Close()
}
