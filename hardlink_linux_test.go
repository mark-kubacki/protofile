// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protofile

import (
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"testing"

	"golang.org/x/sys/unix"
)

func TestHardlink(t *testing.T) {
	// Setup and null hypothesis
	uniqueTempDir := t.TempDir()
	allFilePattern := path.Join(uniqueTempDir, "*")
	var filesInDirectory int
	if m, err := filepath.Glob(allFilePattern); err == nil {
		filesInDirectory = len(m)
	} else {
		t.Fatal("fs.Glob: ", err)
	}

	indexFileName := filepath.Join(uniqueTempDir, "reference")
	indexFile, err := os.Create(indexFileName)
	if err != nil {
		t.Fatalf("os.Create() error == %v", err)
	}
	defer indexFile.Close()
	indexFile.Write([]byte("reference content for comparison"))

	linkedFileName := filepath.Join(uniqueTempDir, "HardLinked")
	// Derive file from indexFile's file descriptor to guarantee
	// our HardLink works not using its name.
	err = Hardlink(os.NewFile(indexFile.Fd(), "/name-does-not-exist"), linkedFileName)
	if err != nil {
		t.Fatalf("Hardlink() error == %v", err)
	}
	if matches, _ := filepath.Glob(allFilePattern); len(matches) != filesInDirectory+2 {
		t.Errorf("Hardlink() didn't materialize the file.\n -- %v", matches)
	}

	if (atomic.LoadInt32(&capState) & unix.CAP_DAC_READ_SEARCH) == 0 {
		t.Skip("Above ran without CAP_DAC_READ_SEARCH.")
		return
	}
	oldCap := atomic.LoadInt32(&capState)
	clearedCap := oldCap &^ unix.CAP_DAC_READ_SEARCH
	for !atomic.CompareAndSwapInt32(&capState, oldCap, clearedCap) {
		oldCap = atomic.LoadInt32(&capState)
		clearedCap = oldCap &^ unix.CAP_DAC_READ_SEARCH
	}

	linkedFileName = filepath.Join(uniqueTempDir, "HardLinked-2")
	err = Hardlink(os.NewFile(indexFile.Fd(), "/name-does-not-exist"), linkedFileName)
	if err != nil {
		t.Fatalf("2: Hardlink() error == %v", err)
	}
	if matches, _ := filepath.Glob(allFilePattern); len(matches) != filesInDirectory+3 {
		t.Errorf("2: Hardlink() with cleared 'capState' didn't materialize the file.\n -- %v", matches)
	}
}
