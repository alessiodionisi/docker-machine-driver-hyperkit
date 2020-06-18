package driver

import (
	"errors"
	"github.com/adnsio/docker-machine-driver-hyperkit/internal/isoinfo"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/mitchellh/go-ps"
)

var (
	errHyperKitProcessNotFound = errors.New("process not found")
)

func (d *Driver) fixFilePermissions(file string) error {
	return os.Chown(file, syscall.Getuid(), syscall.Getegid())
}

func (d *Driver) extractKernelAndRAMDiskFromB2DISO() error {
	isoFile := d.ResolveStorePath("boot2docker.iso")

	isoinfo.ISOInfoBinary = d.ISOInfoBinary

	if err := isoinfo.Extract(isoFile, "/boot/vmlinuz", d.ResolveStorePath("vmlinuz")); err != nil {
		return err
	}

	if err := isoinfo.Extract(isoFile, "/boot/initrd.img", d.ResolveStorePath("initrd.img")); err != nil {
		return err
	}

	return nil
}

func (d *Driver) findHyperKitProcess() (*os.Process, error) {
	pidFile, err := os.Open(d.ResolveStorePath("hyperkit.pid"))
	if err != nil {
		return nil, err
	}
	defer pidFile.Close()

	pidFileData, err := ioutil.ReadAll(pidFile)
	if err != nil {
		return nil, err
	}

	pid, err := strconv.Atoi(string(pidFileData))
	if err != nil {
		return nil, err
	}

	process, err := ps.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	if process == nil {
		return nil, errHyperKitProcessNotFound
	}

	return os.FindProcess(pid)
}
