// +build darwin

package hyperkit

import (
	"testing"
)

func Test_portExtraction(t *testing.T) {
	tests := []struct {
		name    string
		ports   []string
		want    []int
		wantErr error
	}{
		{
			"valid_empty",
			[]string{},
			[]int{},
			nil,
		},
		{
			"valid_list",
			[]string{"10", "20", "30"},
			[]int{10, 20, 30},
			nil,
		},
		{
			"invalid",
			[]string{"8080", "not_an_integer"},
			nil,
			InvalidPortNumberError("not_an_integer"),
		},
	}

	for _, tt := range tests {
		d := NewDriver("", "")
		d.VSockPorts = tt.ports
		got, gotErr := d.extractVSockPorts()
		if !testEq(got, tt.want) {
			t.Errorf("extractVSockPorts() got: %v, want: %v", got, tt.want)
		}
		if gotErr != tt.wantErr {
			t.Errorf("extractVSockPorts() gotErr: %s, wantErr: %s", gotErr.Error(), tt.wantErr.Error())
		}
	}
}

func testEq(a, b []int) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
