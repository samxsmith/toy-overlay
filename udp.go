package main

import (
	"fmt"
	"io"
	"net"
)

func makeUDPRequest(ip net.IP, b []byte) {
	localAddr := &net.UDPAddr{}
	destAddr := &net.UDPAddr{IP: ip, Port: udpPort}
	conn, err := net.DialUDP("udp", localAddr, destAddr)
	if err != nil {
		fmt.Println("error dialing UDP: ", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write(b)
	if err != nil {
		fmt.Println("error sending UDP req: ", err)
	}
}

func startUDPServer(port int) io.ReadCloser {
	portStr := fmt.Sprintf(":%d", port)
	listenAddr, err := net.ResolveUDPAddr("udp", portStr)
	if err != nil {
		pauseOnError(err, "resolving udp address")
	}
	server, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		pauseOnError(err, "starting udp server")
	}
	fmt.Println("Running udp server..")
	return server
}
