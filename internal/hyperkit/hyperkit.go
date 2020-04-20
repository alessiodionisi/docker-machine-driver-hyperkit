package hyperkit

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// HyperKitConfiguration is the HyperKit command configuration
type HyperKitConfiguration struct {
	UUID       string
	CPUs       int
	MemorySize int
	Disk       string
	MACAddress string
	Kernel     string
	RamDisk    string
	Hostname   string
}

var (
	// HyperKitBinary is the path to HyperKit binary
	HyperKitBinary = "hyperkit"
)

// Command returns the configured HyperKit command
func Command(h *HyperKitConfiguration) (*exec.Cmd, error) {
	binaryPath, err := exec.LookPath(HyperKitBinary)
	if err != nil {
		return nil, err
	}

	cmdLine := fmt.Sprintf("loglevel=3 console=ttyS0 console=tty0 host=%s base norestore noembed", h.Hostname)

	args := []string{
		"-A", // ACPI tables
		"-u", // RTC as UTC
		"-H", // Yield vCPU on HLT
		"-P", // Exit vCPU on PAUSE
	}

	args = append(args, "-U", h.UUID)

	args = append(args, "-c", strconv.Itoa(h.CPUs))
	args = append(args, "-m", fmt.Sprintf("%dM", h.MemorySize))

	args = append(args, "-l", "com1,stdio") // LPC device

	args = append(args, "-s", "0,hostbridge") // PCI Host bridge
	args = append(args, "-s", "1,lpc")        // PCI LPC bridge
	args = append(args, "-s", "2,virtio-rnd") // PCI RNG interface

	// PCI network interface
	args = append(args, "-s", fmt.Sprintf("3,virtio-net,,mac=%s", h.MACAddress))

	// PCI block storage interface
	args = append(args, "-s", fmt.Sprintf("4,virtio-blk,%s", h.Disk))

	// Firmware
	args = append(args, "-f", fmt.Sprintf("kexec,%s,%s,earlyprintk=serial %s",
		h.Kernel, h.RamDisk, cmdLine))

	cmd := exec.Command(binaryPath, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
