package telemetry

import (
	"fmt"
	"os"
	"time"
)

type ConverterInterface interface {
	Convert(now time.Time, data map[string]TelemetryData, keys []string)
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
	dataRaw     []byte
	Data        float32
	StartOffset int
	EndOffset   int
}

// DisplayLog Check if flag was passed
func DisplayLog(flagName string, logText string) {
	if os.Getenv("DEBUG_MODE") == flagName {
		fmt.Println(logText)
	}
}
