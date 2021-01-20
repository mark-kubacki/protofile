package protofile

import (
	"os"
	"path"
	"path/filepath"
	"syscall"
	"testing"
)

func TestCreateTemp(t *testing.T) {
	// Setup and null hypothesis
	uniqueTempDir := t.TempDir()
	allFilePattern := path.Join(uniqueTempDir, "*")
	var filesInDirectory int
	if m, err := filepath.Glob(allFilePattern); err == nil {
		filesInDirectory = len(m)
	} else {
		t.Fatal("fs.Glob: ", err)
	}

	// Catch errors on creation, such as yielded by wrong syscall arguments.
	tempFile, err := CreateTemp(uniqueTempDir)
	switch {
	case err == nil: // nop
	case IsTempfileNotSupported(err):
		t.Skipf("Cannot test on this filesystem-kernel combination.\n -- %v", err)
	default:
		t.Errorf("CreateTemp() error = %v, want nil", err)
		return
	}
	if IsTempfileNotSupported(err) != false ||
		IsTempfileNotSupported(syscall.EOPNOTSUPP) != true {
		t.Error("IsTempfileNotSupported() has not been implemented properly")
	}

	// Guarantee no visible effect on the namespace (the directory).
	if matches, _ := filepath.Glob(allFilePattern); len(matches) != filesInDirectory {
		t.Errorf("CreateTemp() should not add to namespace, but did")
	}
	linkedFileName := path.Join(uniqueTempDir, "realname")
	if err := Hardlink(tempFile, linkedFileName); err != nil {
		t.Logf("Hardlink() error = %v", err)
	} else {
		if matches, _ := filepath.Glob(allFilePattern); len(matches) != filesInDirectory+1 {
			t.Errorf("Hardlink() didn't materialize the file.\n -- %v", matches)
		} else {
			os.Remove(linkedFileName)
		}
	}

	// Closing shouldn't have any effect, but be thorough.
	tempFile.Write([]byte("Content that will be discarded."))
	err = tempFile.Close()
	if err != nil {
		t.Errorf("tempFile.Close() error = %v", err)
	}
	if matches, _ := filepath.Glob(allFilePattern); len(matches) != filesInDirectory {
		t.Errorf("Closing tempFile should not leave any entries in the directory, but did")
	}
}
