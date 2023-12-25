package main

import (
	"github.com/bluemanos/simracing-telemetry/src/pkg/converter"
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/bluemanos/simracing-telemetry/src/telemetry/fms2023"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"strconv"
)

func main() {
	debugMode := os.Getenv("DEBUG_MODE")

	fm := &fms2023.ForzaMotorsportHandler{
		TelemetryHandler: telemetry.TelemetryHandler{},
		DebugMode:        debugMode,
	}

	forzaMotorsportPort := getIntPort("TMD_FORZAM")
	if forzaMotorsportPort != -1 {
		err := fm.InitAndRun(forzaMotorsportPort, converter.SetupAdapter(enums.Games.ForzaMotorsport2023()))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func getIntPort(envKey string) int {
	portEnv := os.Getenv(envKey)
	if portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatalln(err)
		}

		return port
	}

	return -1
}
