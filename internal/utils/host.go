package utils

import (
	"log"
	"net"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "192.168.4.159:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
