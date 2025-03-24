package converter

import (
	"fmt"
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
	Clients []*UdpClient
}

type UdpClient struct {
	host       string
	port       int
	addr       *net.UDPAddr
	connection *net.UDPConn
}

func NewUdpForwarder(game enums.Game, adapterConfiguration []string) (*UdpForwarder, error) {
	udpClients := strings.Split(strings.Join(adapterConfiguration[1:], ":"), "&")
	var udpClientsList []*UdpClient
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

		fmt.Printf("UDP Client: %s\n", udpClient)
		udpClientsList = append(udpClientsList, &UdpClient{
			host: udpClientConfiguration[0],
			port: port,
		})
	}

	return &UdpForwarder{
		ConverterData: ConverterData{GameName: game},
		Clients:       udpClientsList,
	}, nil
}

func (udp *UdpForwarder) ChannelInit(now time.Time, channel chan telemetry.GameData, port int) {
	log.Println("UdpForwarder ChannelInit")
	//nolint:gosimple // loop is needed to keep the channel open
	for {
		select {
		case data := <-channel:
			udp.Convert(now, data, port)
		}
	}
}

// Convert converts the data to the UDP clients
func (udp *UdpForwarder) Convert(_ time.Time, data telemetry.GameData, _ int) {
	log.Println("UdpForwarder Convert")
	for _, client := range udp.Clients {
		if client.connection == nil {
			udp.connectToClient(client)
		}
		_, err := client.connection.Write(data.RawData)
		if err != nil {
			if client.connection != nil {
				client.connection.Close()
			}
			client.connection = nil
		}
	}
}

func (udp *UdpForwarder) connectToClient(client *UdpClient) {
	log.Println("UdpForwarder connectToClient")
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
