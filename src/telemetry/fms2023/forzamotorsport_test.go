package fms2023_test

import (
	"log"
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
	fm.Telemetries, fm.Keys = telemetry.Telemetries()

	lines, err := telemetry.ReadLines("fms2023/" + fms2023.DataFormatFile)
	if err != nil {
		log.Fatalf("Error reading format file: %s", err)
	}

	assert.Equal(t, len(fm.Telemetries), len(fm.Keys))
	assert.Equal(t, len(lines), len(fm.Telemetries))
}
