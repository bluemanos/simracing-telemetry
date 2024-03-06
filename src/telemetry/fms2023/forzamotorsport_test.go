package fms2023

import (
	"bufio"
	"encoding/base64"
	"log"
	"os"
	"testing"

	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	_ "github.com/bluemanos/simracing-telemetry/test"
	"github.com/stretchr/testify/assert"
)

func TestInitTelemetries(t *testing.T) {
	fm := &ForzaMotorsportHandler{
		TelemetryHandler: telemetry.TelemetryHandler{},
	}
	fm.Telemetries, fm.Keys = fm.initTelemetry()

	lines, err := telemetry.ReadLines("fms2023/" + dataFormatFile)
	if err != nil {
		log.Fatalf("Error reading format file: %s", err)
	}

	assert.Equal(t, len(fm.Telemetries), len(fm.Keys))
	assert.Equal(t, len(lines), len(fm.Telemetries))
}

func TestProcessBuffer(t *testing.T) {
	fm := &ForzaMotorsportHandler{
		TelemetryHandler: telemetry.TelemetryHandler{},
	}
	fm.Telemetries, fm.Keys = fm.initTelemetry()

	file, err := os.Open("src/telemetry/fms2023/forzamotorsport.udp.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineDecode, _ := base64.StdEncoding.DecodeString(scanner.Text())
		assert.NotPanics(t, func() { fm.processBuffer(lineDecode, 1234) })
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
