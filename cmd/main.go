package main

import (
	"github.com/adnsio/docker-machine-driver-hyperkit/internal/driver"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(driver.New("", ""))
}
