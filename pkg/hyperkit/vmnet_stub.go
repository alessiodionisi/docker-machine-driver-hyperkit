// +build darwin,!cgo

package hyperkit

import (
	"errors"
)

func GetMACAddressFromUUID(UUID string) (string, error) {
	return "", errors.New("Function not supported on CGO_ENABLED=0 binaries")
}
