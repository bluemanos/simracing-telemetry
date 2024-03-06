package server

import "testing"

func TestUdpServer(t *testing.T) {
	t.Run("should return error when could not resolve UDP addr", func(t *testing.T) {
		udpServer := &UDPServer{
			Addr: "invalid",
		}

		err := udpServer.Run(func([]byte, int) {}, 1234)

		if err == nil {
			t.Errorf("Run() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("should run successful", func(t *testing.T) {
		udpServer := &UDPServer{
			Addr: ":1234",
		}

		go func() {
			err := udpServer.Run(func([]byte, int) {}, 1234)
			defer udpServer.Close()

			if err != nil {
				t.Errorf("Run() error = %v", err)
			}
		}()
	})
}
