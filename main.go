package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/songgao/water/waterutil"
)

// pauseOnError pauses the container so that the logs can be read before it is destroyed
func pauseOnError(err error, context string) {
	if err != nil {
		fmt.Printf("Error: %s -> %s\n", context, err)
		time.Sleep(time.Hour * 10000)
	}
}

const (
	flannelRunFile      = "/run/flannel/subnet.env"
	mtu                 = 1300
	tunName             = "tun0"
	udpPort             = 8285
	udpPacketSize       = 1024
	fileWritePermission = 0644
)

var (
	hostname                        = os.Getenv("HOSTNAME")
	nodeCIDRMaskSize                = ""
	dataProvider     dataProviderIf = nil
)

func main() {
	dataProvider = newNodeDataProvider()
	hostNode := dataProvider.getNodeByName(hostname)
	nodeCIDRMaskSize = getCIDRMaskSize(hostNode.podCIDR)
	overlayNetworkCIDR := reduceCIDRSpecificity(hostNode.podCIDR)

	writeCNIPluginFile(overlayNetworkCIDR, hostNode.podCIDR)

	// examples
	// hostNode.podCIDR = 10.244.2.0/24
	// overlayNetworkCIDR & tun addr = 10.244.2.0/16
	// tunRoutingCIDR = 10.244.0.0/16

	tunRoutingCIDR := applyCIDRMask(overlayNetworkCIDR) + "/" + getCIDRMaskSize(overlayNetworkCIDR)
	tun := NewTun(tunName, overlayNetworkCIDR, tunRoutingCIDR, mtu)
	udpServer := startUDPServer(udpPort)

	go listener(tun, mtu, outboundPacketHandler)
	go listener(udpServer, udpPacketSize, func(n int, b []byte) {
		fmt.Println("Inbound packet")
		printIPPacket(b)
		_, err := tun.Write(b[:n])
		if err != nil {
			fmt.Println("error writing to tun: ", err)
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c

	udpServer.Close()
	tun.Close()
}

func writeCNIPluginFile(overlayNetworkCIDR, podCIDR string) {
	// used by the flannel CNI plugin
	fileData := fmt.Sprintf("FLANNEL_NETWORK=%s\nFLANNEL_SUBNET=%s\nFLANNEL_MTU=%d\nFLANNEL_IPMASQ=true", overlayNetworkCIDR, podCIDR, mtu)
	err := ioutil.WriteFile(flannelRunFile, []byte(fileData), fileWritePermission)
	pauseOnError(err, "Write subnet.env")
}

func listener(r io.Reader, bufferSize int, dataHandler func(int, []byte)) {
	for {
		b := make([]byte, bufferSize)
		n, err := r.Read(b)
		if err != nil {
			fmt.Println("Unable to read: ", err)
			continue
		}
		if n == 0 {
			continue
		}
		dataHandler(n, b)
	}
}

func outboundPacketHandler(n int, b []byte) {
	if !waterutil.IsIPv4(b) {
		// for now we won't handle IPv6
		return
	}

	fmt.Println("Outbound packet")
	printIPPacket(b)

	destinationIP := waterutil.IPv4Destination(b)

	destinationCIDR := destinationIP.String() + "/" + nodeCIDRMaskSize
	destinationNode := dataProvider.getNodeMatchingPodCIDR(destinationCIDR)
	if len(destinationNode.internalIP) == 0 {
		fmt.Println("Can't find node for destination pod: ", destinationIP)
		return
	}
	makeUDPRequest(destinationNode.internalIP, b)
}
