package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/telemetry/fms2023"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      os.Getenv("APP_ENVIRONMENT"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	debugMode := os.Getenv("DEBUG_MODE")

	fm := fms2023.NewForzaMotorsportHandler(debugMode)

	forzaMotorsportPort := getIntPort(os.Getenv("TMD_FORZAM"))
	if forzaMotorsportPort != -1 {
		err := fm.InitAndRun(forzaMotorsportPort)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func getIntPort(portEnv string) int {
	if portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatalln(err)
		}

		return port
	}
	return -1
}
