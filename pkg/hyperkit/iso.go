package hyperkit

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hooklift/iso9660"
)

// ExtractFile extracts a file from an ISO
func ExtractFile(isoPath, srcPath, destPath string) error {
	iso, err := os.Open(isoPath)
	if err != nil {
		return err
	}
	defer iso.Close()

	r, err := iso9660.NewReader(iso)
	if err != nil {
		return err
	}

	f, err := findFile(r, srcPath)
	if err != nil {
		return err
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, f.Sys().(io.Reader))
	return err
}

func findFile(r *iso9660.Reader, path string) (os.FileInfo, error) {
	// Look through the ISO for a file with a matching path.
	for f, err := r.Next(); err != io.EOF; f, err = r.Next() {
		// For some reason file paths in the ISO sometimes contain a '.' character at the end, so strip that off.
		if strings.TrimSuffix(f.Name(), ".") == path {
			return f, nil
		}
	}
	return nil, fmt.Errorf("unable to find file %s", path)
}
