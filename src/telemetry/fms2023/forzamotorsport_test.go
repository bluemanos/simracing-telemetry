package fms2023_test

import (
	"bufio"
	"encoding/base64"
	"log"
	"os"
	"testing"

	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/bluemanos/simracing-telemetry/src/telemetry/fms2023"
	_ "github.com/bluemanos/simracing-telemetry/test"
	"github.com/stretchr/testify/assert"
)

func TestInitTelemetries(t *testing.T) {
	fm := &fms2023.ForzaMotorsportHandler{
		TelemetryHandler: telemetry.TelemetryHandler{},
	}
	fm.Telemetries, fm.Keys = fm.InitTelemetry()

	lines, err := telemetry.ReadLines("fms2023/" + fms2023.DataFormatFile)
	if err != nil {
		log.Fatalf("Error reading format file: %s", err)
	}

	assert.Equal(t, len(fm.Telemetries), len(fm.Keys))
	assert.Equal(t, len(lines), len(fm.Telemetries))
}

func TestProcessBuffer(t *testing.T) {
	fm := &fms2023.ForzaMotorsportHandler{
		TelemetryHandler: telemetry.TelemetryHandler{},
	}
	fm.Telemetries, fm.Keys = fm.InitTelemetry()

	file, err := os.Open("src/telemetry/fms2023/forzamotorsport.udp.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//nolint:errcheck
		lineDecode, _ := base64.StdEncoding.DecodeString(scanner.Text())
		assert.NotPanics(t, func() { fm.ProcessBuffer(lineDecode, 1234) })
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
