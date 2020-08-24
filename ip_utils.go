package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/songgao/water/waterutil"
)

const (
	iPv4Size = 32
	byteSize = 8
)

// e.g. 192.244.1.34/24 -> 192.244.1.0
func applyCIDRMask(cidr string) string {
	cidrParts := strings.Split(cidr, "/")
	if len(cidrParts) != 2 {
		pauseOnError(fmt.Errorf("%s is not a valid CIDR", cidr), "Applying Mask")
	}
	maskSize, _ := strconv.Atoi(cidrParts[1])
	mask := net.CIDRMask(maskSize, iPv4Size)
	ip := net.ParseIP(cidrParts[0])
	return ip.Mask(mask).String()
}

// e.g. 10.244.2.0/24 -> 10.244.2.0/16
func reduceCIDRSpecificity(cidr string) string {
	cidrParts := strings.Split(cidr, "/")
	if len(cidrParts) != 2 {
		pauseOnError(fmt.Errorf("%s is not a valid CIDR", cidr), "Applying Mask")
	}
	maskSize, _ := strconv.Atoi(cidrParts[1])
	newMaskSize := strconv.Itoa(maskSize - byteSize)
	return cidrParts[0] + "/" + newMaskSize
}

func getCIDRMaskSize(cidr string) string {
	cidrParts := strings.Split(cidr, "/")
	if len(cidrParts) != 2 {
		pauseOnError(fmt.Errorf("%s is not a valid CIDR", cidr), "Applying Mask")
	}
	return cidrParts[1]
}

func printIPPacket(b []byte) {
	fmt.Println(waterutil.IPv4Destination(b))
}
