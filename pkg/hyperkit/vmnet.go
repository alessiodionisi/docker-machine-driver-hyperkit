// +build darwin,cgo

package hyperkit

import (
	vmnet "github.com/zchee/go-vmnet"
)

func GetMACAddressFromUUID(id string) (string, error) {
	return vmnet.GetMACAddressFromUUID(id)
}
