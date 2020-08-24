// +build linux

package main

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/songgao/water"
)

func NewTun(name, tunAddr, tunRoutingCIDR string, mtu int) io.ReadWriteCloser {
	tun := createTun(name)
	setupTun(name, tunAddr, tunRoutingCIDR, mtu)
	return tun
}

func createTun(name string) io.ReadWriteCloser {
	tunConfig := water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name:        name,
			Persist:     false,
			Permissions: nil,
			MultiQueue:  false,
		},
	}

	tun, err := water.New(tunConfig)
	pauseOnError(err, "Unable to create tun device")
	return tun
}

type tunCommand struct {
	commandString string
	errorIsFatal  bool
}

func setupTun(name, tunAddr, tunRoutingCIDR string, mtu int) {
	mtuStr := strconv.Itoa(mtu)
	tunSetupCommands := []tunCommand{
		{"ip link set dev " + name + " mtu " + mtuStr, true},
		{"ip addr add " + tunAddr + " dev " + name, true},
		{"ip link set dev " + name + " up", true},
		{"ip route add " + tunRoutingCIDR + " dev " + name, false},
	}
	for _, cmd := range tunSetupCommands {
		fmt.Println(cmd)
		err := execCommand(cmd.commandString)
		if cmd.errorIsFatal {
			pauseOnError(err, cmd.commandString)
		}
	}

}

func execCommand(cmd string) error {
	cmdArgs := strings.Split(cmd, " ")
	return exec.Command(cmdArgs[0], cmdArgs[1:]...).Run()
}
