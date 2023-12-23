package server

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

const (
	udpPort = ":6250"
)

func init() {
	udp := NewServer(udpPort)

	go func() {
		udp.Run(func(buffer []byte) {})
	}()
}

func TestServer_Running(t *testing.T) {
	t.Parallel()

	servers := []struct {
		protocol    string
		addr        string
		errExpected error
	}{
		{"tcp", ":1123", errors.New("dial tcp :1123: connect: connection refused")},
		{"udp", udpPort, nil},
	}

	for _, serv := range servers {
		conn, err := net.DialTimeout(serv.protocol, serv.addr, time.Second)
		if err != nil {
			assert.Error(t, serv.errExpected, err)
			continue
		}
		defer conn.Close()
	}
}

func TestServer_Request(t *testing.T) {
	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		{"Sending a simple request returns result", []byte("hello world\n"), []byte("Request received: hello world")},
		{"Sending another simple request works", []byte("goodbye world\n"), []byte("Request received: goodbye world")},
	}

	for _, tc := range tt {
		t.Run(tc.test, func(t *testing.T) {
			conn, err := net.DialTimeout("udp", udpPort, time.Second)
			if err != nil {
				t.Error("could not connect to server: ", err)
			}
			defer conn.Close()

			if _, err := conn.Write(tc.payload); err != nil {
				t.Error("could not write payload to server:", err)
			}
		})
	}
}
