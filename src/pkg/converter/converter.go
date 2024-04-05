package converter

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/afero"
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
	ChannelInit(now time.Time, channel chan telemetry.GameData, port int)
	Convert(now time.Time, data telemetry.GameData, port int)
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
			config, err := NewCsvConverter(game, adapterConfiguration, afero.NewOsFs())
			if err != nil {
				log.Println(err)
				continue
			}
			converters = append(converters, config)
			log.Printf("[%s] CSV adapter configured", game)
		case "mysql":
			config, err := NewMySQLConverter(game, adapterConfiguration)
			if err != nil {
				log.Println(err)
				continue
			}
			converters = append(converters, config)
			log.Printf("[%s] MySQL adapter configured", game)
		case "mysql_bl":
			config, err := NewMysqlBestLapConverter(game, adapterConfiguration)
			if err != nil {
				log.Println(err)
				continue
			}
			converters = append(converters, config)
			log.Printf("[%s] MySQL BL adapter configured", game)
		case "udp":
			config, err := NewUdpForwarder(game, adapterConfiguration)
			if err != nil {
				log.Println(err)
				continue
			}
			converters = append(converters, config)
			log.Printf("[%s] UDP adapter configured", game)
		}
	}
	return converters
}
