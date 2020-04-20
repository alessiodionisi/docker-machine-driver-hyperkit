package driver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/adnsio/docker-machine-driver-hyperkit/internal/dhcpdleases"
	"github.com/adnsio/docker-machine-driver-hyperkit/internal/hyperkit"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/google/uuid"
)

const (
	defaultHyperKitBinary = "hyperkit"
	defaultISOInfoBinary  = "isoinfo"
	defaultDiskSize       = 20000
	defaultMemorySize     = 1024
	defaultCPUs           = 1
)

// Driver is the Docker Machine driver for HyperKit
type Driver struct {
	*drivers.BaseDriver
	HyperKitBinary string
	ISOInfoBinary  string
	Boot2DockerURL string
	DiskSize       int
	MemorySize     int
	CPUs           int
}

// New returns a new HyperKit driver
func New(machineName, storePath string) *Driver {
	return &Driver{
		HyperKitBinary: defaultHyperKitBinary,
		DiskSize:       defaultDiskSize,
		MemorySize:     defaultMemorySize,
		CPUs:           defaultCPUs,
		BaseDriver: &drivers.BaseDriver{
			MachineName: machineName,
			StorePath:   storePath,
		},
	}
}

// Create a host using the driver's config
func (d *Driver) Create() error {
	// download and copy b2d ISO to StorePath
	b2dUtils := mcnutils.NewB2dUtils(d.StorePath)
	if err := b2dUtils.CopyIsoToMachineDir(d.Boot2DockerURL, d.MachineName); err != nil {
		return err
	}

	// extract vmlinuz and initrd.img from b2d ISO
	log.Info("Extracting Kernel and RAM disk from boot2docker ISO...")
	if err := d.extractKernelAndRAMDiskFromB2DISO(); err != nil {
		return err
	}

	sshKeyPath := d.GetSSHKeyPath()
	// generate ssh keys
	log.Info("Generating SSH keys...")
	if err := ssh.GenerateSSHKey(sshKeyPath); err != nil {
		return err
	}

	// create tar buffer with SSH keys for b2d
	diskImageBuffer, err := mcnutils.MakeDiskImage(fmt.Sprintf("%s.pub", sshKeyPath))
	if err != nil {
		return err
	}

	// create disk file
	diskFilePath := d.ResolveStorePath("disk.img")
	log.Infof("Creating disk %s...", diskFilePath)
	diskFile, err := os.Create(diskFilePath)
	if err != nil {
		return err
	}
	defer diskFile.Close()

	if _, err := diskFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// write tar data to disk file
	if _, err = diskFile.Write(diskImageBuffer.Bytes()); err != nil {
		return err
	}

	// truncate disk file
	if err := os.Truncate(diskFilePath, int64(d.DiskSize*1000000)); err != nil {
		return err
	}

	// fix permissions of store path
	storePath := d.ResolveStorePath("")
	log.Infof("Fixing permissions of %s...", storePath)
	if err := d.fixFilePermissions(storePath); err != nil {
		return err
	}

	// fix permissions of all files
	files, err := ioutil.ReadDir(storePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(storePath, file.Name())
		log.Infof("Fixing permissions of %s...", filePath)
		if err := d.fixFilePermissions(filePath); err != nil {
			return err
		}
	}

	return d.Start()
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "hyperkit"
}

// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:   "hyperkit-binary",
			Usage:  "Path of the HyperKit binary.",
			Value:  defaultHyperKitBinary,
			EnvVar: "HYPERKIT_BINARY",
		},
		mcnflag.StringFlag{
			Name:   "hyperkit-isoinfo-binary",
			Usage:  "Path of the isoinfo binary.",
			Value:  defaultISOInfoBinary,
			EnvVar: "HYPERKIT_ISOINFO_BINARY",
		},
		mcnflag.StringFlag{
			Name:   "hyperkit-boot2docker-url",
			Usage:  "URL of the boot2docker ISO. Defaults to the latest version.",
			EnvVar: "HYPERKIT_BOOT2DOCKER_URL",
		},
		mcnflag.IntFlag{
			Name:   "hyperkit-disk-size",
			Usage:  "Guest disk size in megabytes.",
			Value:  defaultDiskSize,
			EnvVar: "HYPERKIT_DISK_SIZE",
		},
		mcnflag.IntFlag{
			Name:   "hyperkit-memory-size",
			Usage:  "Guest physical memory size in megabytes.",
			Value:  defaultMemorySize,
			EnvVar: "HYPERKIT_MEMORY_SIZE",
		},
		mcnflag.IntFlag{
			Name:   "hyperkit-cpus",
			Usage:  "Number of guest virtual CPUs. The default is 1 and the maximum is 16.",
			Value:  defaultCPUs,
			EnvVar: "HYPERKIT_CPUS",
		},
	}
}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetURL returns a Docker compatible host URL for connecting to this host
// e.g. tcp://1.2.3.4:2376
func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	if ip == "" {
		return "", nil
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

// GetState returns the state that the host is in (running, stopped, etc)
func (d *Driver) GetState() (state.State, error) {
	_, err := d.findHyperKitProcess()
	if err != nil {
		if err == errHyperKitProcessNotFound {
			return state.Stopped, nil
		}

		return state.Error, err
	}

	return state.Running, nil
}

// Kill stops a host forcefully
func (d *Driver) Kill() error {
	return d.Stop()
}

// Remove a host
func (d *Driver) Remove() error {
	return d.Stop()
}

// Restart a host. This may just call Stop(); Start() if the provider does not
// have any special restart behaviour.
func (d *Driver) Restart() error {
	if err := d.Stop(); err != nil {
		return err
	}

	return d.Start()
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.HyperKitBinary = opts.String("hyperkit-binary")
	d.ISOInfoBinary = opts.String("hyperkit-isoinfo-binary")
	d.Boot2DockerURL = opts.String("hyperkit-boot2docker-url")
	d.DiskSize = opts.Int("hyperkit-disk-size")
	d.MemorySize = opts.Int("hyperkit-memory-size")
	d.CPUs = opts.Int("hyperkit-cpus")
	d.SSHUser = "docker"
	d.SetSwarmConfigFromFlags(opts)

	return nil
}

// Start a host
func (d *Driver) Start() error {
	machineUUID := uuid.NewSHA1(uuid.Nil, []byte(d.GetMachineName()))
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		machineUUID[0], machineUUID[1], machineUUID[2],
		machineUUID[3], machineUUID[4], machineUUID[5])

	log.Infof("UUID: %s", machineUUID)
	log.Infof("MAC Address: %s", macAddress)

	disk := d.ResolveStorePath("disk.img")
	kernel := d.ResolveStorePath("vmlinuz")
	ramDisk := d.ResolveStorePath("initrd.img")

	hyperkit.HyperKitBinary = d.HyperKitBinary
	hkConfig := &hyperkit.HyperKitConfiguration{
		UUID:       machineUUID.String(),
		CPUs:       d.CPUs,
		MemorySize: d.MemorySize,
		Disk:       disk,
		MACAddress: macAddress,
		Kernel:     kernel,
		RamDisk:    ramDisk,
		Hostname:   d.GetMachineName(),
	}

	hkCMD, err := hyperkit.Command(hkConfig)
	if err != nil {
		return err
	}
	log.Infof("CMD: %s", hkCMD.String())

	if err := hkCMD.Start(); err != nil {
		return err
	}

	log.Infof("PID: %d", hkCMD.Process.Pid)

	pidFile, err := os.Create(d.ResolveStorePath("hyperkit.pid"))
	if err != nil {
		return err
	}
	defer pidFile.Close()

	if _, err := pidFile.Write([]byte(strconv.Itoa(hkCMD.Process.Pid))); err != nil {
		return err
	}

	for i := 0; i < 30; i++ {
		log.Infof("Finding IP Address... (%d/30)", i+1)
		d.IPAddress, err = dhcpdleases.FindIPAddress(dhcpdleases.DefaultDHCPDLeasesFile, macAddress)
		if err == nil {
			break
		}
		if err != dhcpdleases.ErrLeaseNotFound {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return err
	}

	log.Infof("IP Address: %s", d.IPAddress)

	return nil
}

// Stop a host gracefully
func (d *Driver) Stop() error {
	process, err := d.findHyperKitProcess()
	if err != nil {
		if err == errHyperKitProcessNotFound {
			return nil
		}
		return err
	}

	if err := process.Kill(); err != nil {
		return err
	}

	return nil
}
