// +build darwin

package main

import (
	"fmt"
	"os"

	"github.com/adnsio/docker-machine-driver-hyperkit/pkg/hyperkit"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("version:", hyperkit.GetVersion())
		fmt.Println("commit:", hyperkit.GetGitCommitID())
		return
	}

	plugin.RegisterDriver(hyperkit.NewDriver("", ""))
}
