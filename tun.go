// +build !linux

package main

import (
	"errors"
	"io"
)

var errNotImplemented = errors.New("Not implemented")

func NewTun(name, tunAddr, tunRoutingCIDR string, mtu int) (r io.ReadWriteCloser) {
	pauseOnError(errNotImplemented, "NewTun")
	return
}
