package converter

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/pkg/errors"
)

var ErrInvalidUDPAdapterConfiguration = errors.New("[UDP] invalid adapter configuration")

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

func NewUdpForwarder(game enums.Game, adapterConfiguration []string) (*UdpForwarder, error) {
	udpClients := strings.Split(strings.Join(adapterConfiguration[1:], ":"), "&")
	var udpClientsList []UdpClient
	for _, udpClient := range udpClients {
		udpClientConfiguration := strings.Split(udpClient, ":")
		if len(udpClientConfiguration) != 2 {
			return nil, errors.Wrapf(ErrInvalidUDPAdapterConfiguration,
				"[%s] Wrong UDP adapter configuration: %s", game, udpClient,
			)
		}
		port, err := strconv.Atoi(udpClientConfiguration[1])
		if err != nil {
			return nil, errors.Wrapf(ErrInvalidUDPAdapterConfiguration,
				"[%s] Wrong UDP adapter configuration: %s", game, udpClient,
			)
		}

		udpClientsList = append(udpClientsList, UdpClient{
			host: udpClientConfiguration[0],
			port: port,
		})
	}

	return &UdpForwarder{
		ConverterData: ConverterData{GameName: game},
		Clients:       udpClientsList,
	}, nil
}

// Convert converts the data to the UDP clients
func (udp *UdpForwarder) Convert(_ time.Time, data telemetry.GameData, _ int) {
	for _, client := range udp.Clients {
		client := client
		if client.connection == nil {
			udp.connectToClient(&client)
			defer client.connection.Close()
		}
		_, err := client.connection.Write(data.RawData)
		if err != nil {
			return
		}
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
