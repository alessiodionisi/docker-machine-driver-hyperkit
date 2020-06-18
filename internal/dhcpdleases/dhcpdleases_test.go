package dhcpdleases

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var (
	dhcpdLeasesFile = []byte(`{
	name=vm1
	ip_address=1.2.3.4
	hw_address=1,a1:b2:c3:d4:e5:f6
	identifier=1,a1:b2:c3:d4:e5:f6
	lease=0x00000000
}
{
	name=vm2
	ip_address=2.2.3.4
	hw_address=1,a2:b2:c3:d4:e5:f6
	identifier=1,a2:b2:c3:d4:e5:f6
	lease=0x00000000
}`)
)

func TestFindIPAddress(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "*.dhcpd_leases")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tmpFile: %s", tmpFile.Name())
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(dhcpdLeasesFile); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name   string
		Input  string
		Output string
		Error  error
	}{
		{
			Name:   "Found1",
			Input:  "a1:b2:c3:d4:e5:f6",
			Output: "1.2.3.4",
		},
		{
			Name:   "Found2",
			Input:  "a2:b2:c3:d4:e5:f6",
			Output: "2.2.3.4",
		},
		{
			Name:   "NotFound",
			Input:  "00:00:00:00:00:00",
			Output: "",
			Error:  ErrLeaseNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			out, err := FindIPAddress(tmpFile.Name(), tt.Input)
			if !reflect.DeepEqual(err, tt.Error) {
				t.Errorf("FindIPAddress() error:\ngot  %v\nwant %v", err, tt.Error)
			} else if !reflect.DeepEqual(out, tt.Output) {
				t.Errorf("FindIPAddress() output:\ngot  %q\nwant %q", out, tt.Output)
			}
		})
	}
}
