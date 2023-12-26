package converter

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/spf13/afero"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
)

var gameEnvKeys = map[enums.Game]gameConfiguration{
	enums.Games.ForzaMotorsport2023(): {
		AdaptersEnvKey: "TMD_FORZAM_ADAPTERS",
		DatabaseTable:  "tmd_forzamotorsport2023",
	},
}

type gameConfiguration struct {
	AdaptersEnvKey string
	DatabaseTable  string
}

type ConverterInterface interface {
	Convert(now time.Time, data map[string]telemetry.TelemetryData, keys []string)
}

type ConverterData struct {
	GameName enums.Game
}

// SetupAdapter sets up game adapters like CSV export, MySQL export, etc.
func SetupAdapter(game enums.Game) []telemetry.ConverterInterface {
	adapters := strings.Split(os.Getenv(gameEnvKeys[game].AdaptersEnvKey), ",")

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
				},
				Fs:        afero.NewOsFs(),
				FilePath:  adapterConfiguration[1],
				Retention: enums.RetentionType(adapterConfiguration[2]),
			})
			log.Printf("[%s] CSV adapter configured", game)
		case "mysql":
			if len(adapterConfiguration) != 6 {
				log.Printf("[%s] Wrong MySQL adapter configuration", game)
				continue
			}

			var db *sql.DB
			if flag.Lookup("test.v") == nil {
				var err error
				db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", adapterConfiguration[1], adapterConfiguration[2], adapterConfiguration[3], adapterConfiguration[4], adapterConfiguration[5]))
				if err != nil {
					log.Println(err)
					continue
				}
			}

			converters = append(converters, MySqlConverter{
				ConverterData: ConverterData{
					GameName: game,
				},
				User:      adapterConfiguration[1],
				Password:  adapterConfiguration[2],
				Host:      adapterConfiguration[3],
				Port:      adapterConfiguration[4],
				Database:  adapterConfiguration[5],
				TableName: gameEnvKeys[game].DatabaseTable,
				connector: db,
			})
			log.Printf("[%s] MySQL adapter configured", game)
		}
	}
	return converters
}
