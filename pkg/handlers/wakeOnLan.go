// wakeOnLan.go
package handlers

import (
	"net"
)

func SendWol(macAddr string) error {

	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return err
	}

	var packet []byte

	packet = append(packet, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}...)

	for range 16 {
		packet = append(packet, hwAddr...)
	}

	addr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(packet)
	if err != nil {
		return err
	}

	return nil

}
