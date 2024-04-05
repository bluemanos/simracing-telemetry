package fms2023

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/converter"
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/pkg/server"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
)

const DataFormatFile = "forzamotorsport"

type ForzaMotorsportHandler struct {
	telemetry.TelemetryHandler
	DebugMode string
}

var gameTelemetryData = make(chan telemetry.GameData)

// NewForzaMotorsportHandler creates a new ForzaMotorsportHandler
func NewForzaMotorsportHandler(debugMode string) *ForzaMotorsportHandler {
	return &ForzaMotorsportHandler{
		TelemetryHandler: telemetry.TelemetryHandler{
			Adapters: converter.SetupAdapter(enums.Games.ForzaMotorsport2023()),
		},
		DebugMode: debugMode,
	}
}

// InitAndRun starts the ForzaMotorsportHandler
func (fm *ForzaMotorsportHandler) InitAndRun(port int) error {
	udpServer := server.NewServer("0.0.0.0:" + strconv.Itoa(port))
	fm.Telemetries, fm.Keys = fm.InitTelemetry()

	log.Printf("Forza data out server listening on %s:%d, waiting for Forza data...\n", telemetry.GetOutboundIP(), port)

	err := udpServer.Run(fm.ProcessChannel, port)
	defer udpServer.Close()
	if err != nil {
		return err
	}
	return nil
}

// InitTelemetry initializes the telemetry data
func (fm *ForzaMotorsportHandler) InitTelemetry() (map[string]telemetry.TelemetryData, []string) {
	lines, err := telemetry.ReadLines("fms2023/" + DataFormatFile)
	if err != nil {
		log.Fatalf("Error reading format file: %s", err)
	}

	telemetryArray := make(map[string]telemetry.TelemetryData, len(lines))
	telemetryKeys := make([]string, len(lines))
	startOffset := 0
	endOffset := 0
	dataLength := 0

	for i, line := range lines {
		dataFormat := strings.Split(line, " ")
		dataType := dataFormat[0]
		dataName := dataFormat[1]

		switch dataType {
		case "S32", "U32", "F32":
			dataLength = 4
		case "U16":
			dataLength = 2
		case "U8", "S8":
			dataLength = 1
		default:
			log.Fatalf("ForzaMotorsportHandler Error: Unknown data type: %s\n", dataType)
		}
		endOffset = endOffset + dataLength
		startOffset = endOffset - dataLength

		telemItem := telemetry.TelemetryData{
			Position:    i,
			Name:        dataName,
			DataType:    dataType,
			StartOffset: startOffset,
			EndOffset:   endOffset,
		}
		telemetryArray[dataName] = telemItem
		telemetryKeys[i] = dataName
	}

	return telemetryArray, telemetryKeys
}

func (fm *ForzaMotorsportHandler) ProcessChannel(channel chan []byte, port int) {
	fmt.Println("ForzaMotorsportHandler ProcessChannel")
	for _, adapter := range fm.Adapters {
		go adapter.ChannelInit(time.Now(), gameTelemetryData, port)
	}

	//nolint:gosimple // loop is needed to keep the channel open
	for {
		select {
		case data := <-channel:
			fm.ProcessBuffer(data, port)
		}
	}
}

// ProcessBuffer processes the received data
func (fm *ForzaMotorsportHandler) ProcessBuffer(buffer []byte, port int) {
	tempTelemetry := make(map[string]float32, len(fm.Telemetries))

	for i, telemetryObj := range fm.Telemetries {
		data := buffer[telemetryObj.StartOffset:telemetryObj.EndOffset]

		var value float32
		switch telemetryObj.DataType {
		case "F32":
			value = math.Float32frombits(binary.LittleEndian.Uint32(data))
		case "U8":
			value = float32(data[0])
		case "S8":
			value = float32(int8(data[0]))
		case "U16":
			value = float32(binary.LittleEndian.Uint16(data))
		default:
			value = float32(binary.LittleEndian.Uint32(data))
		}

		tempTelemetry[i] = value
	}

	if tempTelemetry["IsRaceOn"] == 0 {
		return
	}

	data := telemetry.GameData{
		Keys:    fm.Keys,
		Data:    tempTelemetry,
		RawData: buffer,
	}
	gameTelemetryData <- data
}
