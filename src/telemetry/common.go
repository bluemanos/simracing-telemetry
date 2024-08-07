package telemetry

import (
	"fmt"
	"os"
	"time"
)

type GameData struct {
	Keys    []string
	Data    map[string]float32
	RawData []byte
}

type ConverterInterface interface {
	ChannelInit(now time.Time, channel chan GameData, port int)
	Convert(now time.Time, data GameData, port int)
}

type TelemetryInterface interface {
	InitAndRun(port int) error
}

type TelemetryHandler struct {
	Telemetries map[string]TelemetryData
	Keys        []string
	Adapters    []ConverterInterface
}

type TelemetryData struct {
	Position    int
	Name        string
	DataType    string
	StartOffset int
	EndOffset   int
}

// DisplayLog Check if flag was passed
func DisplayLog(flagName string, logText any) {
	if os.Getenv("DEBUG_MODE") == flagName {
		fmt.Println(logText)
	}
}
