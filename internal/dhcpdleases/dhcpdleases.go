package dhcpdleases

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

const (
	// DefaultDHCPDLeasesFile is the default macOS dhcpd leases file
	DefaultDHCPDLeasesFile = "/var/db/dhcpd_leases"
)

var (
	ErrLeaseNotFound = errors.New("lease not found in dhcpd_leases file")
)

type DHCPLease struct {
	Name       string
	IPAddress  string
	HWAddress  string
	Identifier string
	Lease      string
}

// FindIPAddress returns the IP Address for the provided MAC Address
func FindIPAddress(dhcpdLeasesFile string, macAddress string) (string, error) {
	file, err := os.Open(dhcpdLeasesFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	leadingZerosRegexp := regexp.MustCompile(`0([A-Fa-f0-9](:|$))`)
	macAddressWithoutLeadingZeros := leadingZerosRegexp.ReplaceAllString(macAddress, "$1")

	dhcpLease := &DHCPLease{}
	dhcpLeases := make([]DHCPLease, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "{" {
			dhcpLease = new(DHCPLease)
			continue
		} else if line == "}" {
			dhcpLeases = append(dhcpLeases, *dhcpLease)
			continue
		}

		lineParts := strings.SplitN(line, "=", 2)

		key, value := lineParts[0], lineParts[1]
		switch key {
		case "name":
			dhcpLease.Name = value
		case "ip_address":
			dhcpLease.IPAddress = value
		case "hw_address":
			dhcpLease.HWAddress = value[2:]
		case "identifier":
			dhcpLease.Identifier = value
		case "lease":
			dhcpLease.Lease = value
		}
	}

	if scanner.Err() != nil {
		return "", err
	}

	for _, lease := range dhcpLeases {
		if lease.HWAddress == macAddressWithoutLeadingZeros {
			return lease.IPAddress, nil
		}
	}

	return "", ErrLeaseNotFound
}
