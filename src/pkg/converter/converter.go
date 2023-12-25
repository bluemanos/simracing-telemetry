package converter

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/spf13/afero"
)

var gameEnvKeys = map[enums.Game]string{
	enums.Games.ForzaMotorsport2023(): "TMD_FORZAM_ADAPTERS",
}

type ConverterInterface interface {
	Convert(now time.Time, data map[string]telemetry.TelemetryData, keys []string)
}

type ConverterData struct {
	GameName enums.Game
	Fs       afero.Fs
}

func SetupAdapter(game enums.Game) []telemetry.ConverterInterface {
	adapters := strings.Split(os.Getenv(gameEnvKeys[game]), ",")

	var converters []telemetry.ConverterInterface

	for _, adapter := range adapters {
		adapterConfiguration := strings.Split(adapter, ":")
		switch adapterConfiguration[0] {
		case "csv":
			if len(adapterConfiguration) != 3 {
				log.Printf("[%s] Wrong CSV adapter configuration", game)
				continue
			}

			converters = append(converters, CsvConverter{
				ConverterData: ConverterData{
					GameName: game,
					Fs:       afero.NewOsFs(),
				},
				FilePath:  adapterConfiguration[1],
				Retention: enums.RetentionType(adapterConfiguration[2]),
			})
		}
	}
	return converters
}
