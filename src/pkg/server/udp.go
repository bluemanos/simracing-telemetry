package server

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type UDPServer struct {
	Addr   string
	server *net.UDPConn
}

var udpBuffer = make(chan []byte)

// Run starts the UDP server.
func (u *UDPServer) Run(fn HandleConnection, port int) (err error) {
	laddr, err := net.ResolveUDPAddr("udp", u.Addr)
	if err != nil {
		return errors.New("could not resolve UDP addr")
	}

	u.server, err = net.ListenUDP("udp", laddr)
	if err != nil {
		return errors.New("could not listen on UDP")
	}

	fmt.Println("UPD fn goroutine")
	go fn(udpBuffer, port)

	for {
		buf := make([]byte, 2048)
		n, conn, err := u.server.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			break
		}
		if conn == nil {
			log.Println("UDP: no connection")
			continue
		}

		udpBuffer <- buf[:n]
	}
	return nil
}

// Close ensures that the UDPServer is shut down gracefully.
func (u *UDPServer) Close() error {
	return u.server.Close()
}
