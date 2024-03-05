package converter

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/telemetry"
)

type UdpForwarder struct {
	ConverterData
	Clients []UdpClient
}

type UdpClient struct {
	host       string
	port       int
	addr       *net.UDPAddr
	connection *net.UDPConn
}

// Convert converts the data to the UDP clients
func (udp *UdpForwarder) Convert(_ time.Time, data telemetry.GameData, port int) {
	for _, client := range udp.Clients {
		if client.connection == nil {
			udp.connectToClient(&client)
			defer client.connection.Close()
		}
		client.connection.Write(data.RawData)
	}
}

func (udp *UdpForwarder) connectToClient(client *UdpClient) {
	var err error
	address := client.host + ":" + strconv.Itoa(client.port)
	client.addr, err = net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println(err)
	}

	client.connection, err = net.DialUDP("udp", nil, client.addr)
	if err != nil {
		log.Println(err)
	}
}
