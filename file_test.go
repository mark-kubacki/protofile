// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protofile // import "blitznote.com/src/protofile"

import (
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Generates a new temporary file name without a path.
func tempFileName() string {
	buffer := make([]byte, 16)
	_, _ = rand.Read(buffer)
	for i := range buffer {
		buffer[i] = (buffer[i] % 25) + 97 // a–z
	}
	return string(buffer)
}

func TestGeneralizedProtoFile(t *testing.T) {
	scratchDir := os.TempDir()

	Convey("GeneralizedProtoFile", t, func() {
		Convey("creates a new file", func() {
			filename := tempFileName()
			fp, err := intentNewUniversal(scratchDir, filename)
			defer func() {
				os.Remove(filepath.Join(scratchDir, filename))
			}()
			defer func() {
				os.Remove(filepath.Join(scratchDir, "."+filename))
			}()
			So(err, ShouldBeNil)
			So(fp, ShouldNotBeNil)
			if fp == nil {
				return
			}
			f := *fp

			n, err := io.Copy(f, strings.NewReader("DELME"))
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 5)

			err = f.Persist()
			So(err, ShouldBeNil)
		})

		Convey("the file is not in visible namespace until persisted", func() {
			filename := tempFileName()
			fp, err := intentNewUniversal(scratchDir, filename)
			So(err, ShouldBeNil)
			defer func() {
				os.Remove(filepath.Join(scratchDir, filename))
			}()
			defer func() {
				os.Remove(filepath.Join(scratchDir, "."+filename))
			}()
			if fp == nil {
				So(fp, ShouldNotBeNil)
				return
			}
			f := *fp

			_, err = os.Stat(filepath.Join(scratchDir, filename))
			So(os.IsNotExist(err), ShouldBeTrue)
			io.Copy(f, strings.NewReader("DELME"))
			_, err = os.Stat(filepath.Join(scratchDir, filename))
			So(os.IsNotExist(err), ShouldBeTrue)

			err = f.Persist()
			So(err, ShouldBeNil)
			_, err = os.Stat(filepath.Join(scratchDir, filename))
			So(os.IsNotExist(err), ShouldBeFalse)
		})

		Convey("the file will not materialize after having been zapped", func() {
			filename := tempFileName()
			fp, err := intentNewUniversal(scratchDir, filename)
			So(err, ShouldBeNil)
			defer func() {
				os.Remove(filepath.Join(scratchDir, filename))
			}()
			defer func() {
				os.Remove(filepath.Join(scratchDir, "."+filename))
			}()
			if fp == nil {
				So(fp, ShouldNotBeNil)
				return
			}
			f := *fp

			io.Copy(f, strings.NewReader("DELME"))

			err = f.Zap()
			So(err, ShouldBeNil)
			_, err = os.Stat(filepath.Join(scratchDir, filename))
			So(os.IsNotExist(err), ShouldBeTrue)
		})
	})
}
