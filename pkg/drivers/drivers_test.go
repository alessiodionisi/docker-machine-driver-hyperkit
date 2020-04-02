package drivers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/minikube/pkg/minikube/tests"
)

func Test_createDiskImage(t *testing.T) {
	tmpdir := tests.MakeTempDir()
	defer os.RemoveAll(tmpdir)

	sshPath := filepath.Join(tmpdir, "ssh")
	if err := ioutil.WriteFile(sshPath, []byte("mysshkey"), 0644); err != nil {
		t.Fatalf("writefile: %v", err)
	}
	diskPath := filepath.Join(tmpdir, "disk")

	sizeInMb := 100
	sizeInBytes := int64(sizeInMb) * 1000000
	if err := createRawDiskImage(sshPath, diskPath, sizeInMb); err != nil {
		t.Errorf("createDiskImage() error = %v", err)
	}
	fi, err := os.Lstat(diskPath)
	if err != nil {
		t.Errorf("Lstat() error = %v", err)
	}
	if fi.Size() != sizeInBytes {
		t.Errorf("Disk size is %v, want %v", fi.Size(), sizeInBytes)
	}
}
