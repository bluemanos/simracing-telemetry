package fms2023

import (
	"encoding/binary"
	"log"
	"math"
	"strconv"
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
	fm.TelemetryHandler.Telemetries, fm.TelemetryHandler.Keys = telemetry.Telemetries()

	log.Printf("Forza data out server listening on %s:%d, waiting for Forza data...\n", telemetry.GetOutboundIP(), port)

	err := udpServer.Run(fm.ProcessChannel, port)
	defer udpServer.Close()
	if err != nil {
		return err
	}
	return nil
}

func (fm *ForzaMotorsportHandler) ProcessChannel(channel chan []byte, port int) {
	for _, adapter := range fm.TelemetryHandler.Adapters {
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
	tempTelemetry := make(map[string]float32, len(fm.TelemetryHandler.Telemetries))

	for i, telemetryObj := range fm.TelemetryHandler.Telemetries {
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
		Keys:    fm.TelemetryHandler.Keys,
		Data:    tempTelemetry,
		RawData: buffer,
	}
	gameTelemetryData <- data
}
