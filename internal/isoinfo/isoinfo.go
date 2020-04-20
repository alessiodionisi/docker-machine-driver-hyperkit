package isoinfo

import (
	"os"
	"os/exec"
)

var (
	// ISOInfoBinary is the path to isoinfo binary
	ISOInfoBinary = "isoinfo"
)

// Extract extract files from ISO
func Extract(file string, from string, to string) error {
	binaryPath, err := exec.LookPath(ISOInfoBinary)
	if err != nil {
		return err
	}

	outFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer outFile.Close()

	args := []string{
		"-R",
	}

	args = append(args, "-i", file)
	args = append(args, "-x", from)

	cmd := exec.Command(binaryPath, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
