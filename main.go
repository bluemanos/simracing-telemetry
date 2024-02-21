package main

import (
	"flag"
	"fmt"
	"github.com/bluemanos/simracing-telemetry/test"
	"log"
	"os"
	"strconv"

	"github.com/bluemanos/simracing-telemetry/src/telemetry/fms2023"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if flag.Lookup("test-server") != nil {
		fmt.Printf("%+v\n", "echo server")
		test.Run()
		return
	}

	debugMode := os.Getenv("DEBUG_MODE")

	fm := fms2023.NewForzaMotorsportHandler(debugMode)

	forzaMotorsportPort := getIntPort("TMD_FORZAM")
	if forzaMotorsportPort != -1 {
		err := fm.InitAndRun(forzaMotorsportPort)
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
