package converter_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/bluemanos/simracing-telemetry/src/pkg/converter"
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var testExpectedMissingCsv = "[fms2023] Wrong CSV adapter configuration"

func TestSetupAdapter(t *testing.T) {
	tt := []struct {
		testName         string
		game             enums.Game
		setup            func(t *testing.T)
		expectedAdapters []telemetry.ConverterInterface
		expectedLogs     *string
	}{
		{
			testName: "forza csv adapter success",
			game:     enums.Games.ForzaMotorsport2023(),
			setup: func(t *testing.T) {
				t.Setenv("TMD_FORZAM_ADAPTERS", "csv:/var/log/simracing-telemetry/fms2023-daily.csv:daily")
			},
			expectedAdapters: []telemetry.ConverterInterface{
				&converter.CsvConverter{
					ConverterData: converter.ConverterData{
						GameName: enums.Games.ForzaMotorsport2023(),
					},
					Fs:        afero.NewOsFs(),
					FilePath:  "/var/log/simracing-telemetry/fms2023-daily.csv",
					Retention: enums.RetentionType("daily"),
				},
			},
		},
		{
			testName: "forza csv adapter missing retention",
			game:     enums.Games.ForzaMotorsport2023(),
			setup: func(t *testing.T) {
				t.Setenv("TMD_FORZAM_ADAPTERS", "csv:/var/log/simracing-telemetry/fms2023-daily.csv")
			},
			expectedAdapters: nil,
			expectedLogs:     &testExpectedMissingCsv,
		},
		{
			testName: "forza mysql adapter",
			game:     enums.Games.ForzaMotorsport2023(),
			setup: func(t *testing.T) {
				t.Setenv("TMD_FORZAM_ADAPTERS", "mysql:user:pass:db:3306:app")
			},
			expectedAdapters: []telemetry.ConverterInterface{
				&converter.MySQLConverter{
					ConverterData: converter.ConverterData{
						GameName: enums.Games.ForzaMotorsport2023(),
					},
					User:      "user",
					Password:  "pass",
					Host:      "db",
					Port:      "3306",
					Database:  "app",
					TableName: "tmd_forzamotorsport2023",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			tc.setup(t)
			var buf bytes.Buffer
			if tc.expectedLogs != nil {
				log.SetOutput(&buf)
				defer func() {
					log.SetOutput(os.Stderr)
				}()
			}

			adapters := converter.SetupAdapter(tc.game)
			assert.Equal(t, tc.expectedAdapters, adapters)

			if tc.expectedLogs != nil {
				assert.Contains(t, buf.String(), *tc.expectedLogs)
			}
		})
	}
}
