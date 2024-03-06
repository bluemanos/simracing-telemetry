package converter

import (
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCsvFilename(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	fs.MkdirAll("/var/log/simracing-telemetry", 0755)
	fs.MkdirAll("relative/path/simracing-telemetry", 0755)

	tt := []struct {
		testName         string
		time             time.Time
		converterSetup   CsvConverter
		expectedFilePath string
		expectedError    error
	}{
		{
			testName: "daily: path with trailing slash",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry/",
				Retention: enums.RetentionTypes.Daily(),
			},
			expectedFilePath: "/var/log/simracing-telemetry/fms2023-daily-2021-01-01.csv",
		},
		{
			testName: "daily: path without trailing slash",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry",
				Retention: enums.RetentionTypes.Daily(),
			},
			expectedFilePath: "/var/log/simracing-telemetry/fms2023-daily-2021-01-01.csv",
		},
		{
			testName: "daily: path with filename",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry/test.txt",
				Retention: enums.RetentionTypes.Daily(),
			},
			expectedFilePath: "",
			expectedError:    errInvalidFilePath,
		},
		{
			testName: "daily: path not exists",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry123",
				Retention: enums.RetentionTypes.Daily(),
			},
			expectedFilePath: "",
			expectedError:    errInvalidFilePath,
		},
		{
			testName: "daily: relative path not exists",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "var/log/not-exists/simracing-telemetry",
				Retention: enums.RetentionTypes.Daily(),
			},
			expectedFilePath: "",
			expectedError:    errInvalidFilePath,
		},
		{
			testName: "daily: relative path with trailing slash",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "relative/path/simracing-telemetry/",
				Retention: enums.RetentionTypes.Daily(),
			},
			expectedFilePath: "relative/path/simracing-telemetry/fms2023-daily-2021-01-01.csv",
		},
		{
			testName: "no retention: path with trailing slash",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry/",
				Retention: enums.RetentionTypes.None(),
			},
			expectedFilePath: "/var/log/simracing-telemetry/fms2023.csv",
		},
		{
			testName: "no retention: path with trailing slash",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry",
				Retention: enums.RetentionTypes.None(),
			},
			expectedFilePath: "/var/log/simracing-telemetry/fms2023.csv",
		},
		{
			testName: "no retention: path not exists",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "var/log/not-exists/simracing-telemetry",
				Retention: enums.RetentionTypes.None(),
			},
			expectedFilePath: "",
			expectedError:    errInvalidFilePath,
		},
		{
			testName: "no retention: path with filename",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			converterSetup: CsvConverter{
				ConverterData: ConverterData{
					GameName: enums.Games.ForzaMotorsport2023(),
				},
				Fs:        fs,
				FilePath:  "/var/log/simracing-telemetry/newfile.csv",
				Retention: enums.RetentionTypes.None(),
			},
			expectedFilePath: "/var/log/simracing-telemetry/newfile.csv",
		},
	}

	for i := range tt {
		test := tt[i]
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			filePath, err := test.converterSetup.correctFilePath(test.time)
			assert.Equal(t, test.expectedFilePath, filePath)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestCsvConvert(t *testing.T) {
	now := time.Date(2023, 12, 24, 0, 1, 2, 333, time.UTC)

	fs := afero.NewMemMapFs()
	fs.MkdirAll("/var/www/simracing-telemetry", 0755)

	converter := CsvConverter{
		ConverterData: ConverterData{
			GameName: enums.Games.ForzaMotorsport2023(),
		},
		Fs:        fs,
		FilePath:  "/var/www/simracing-telemetry/test.csv",
		Retention: enums.RetentionTypes.None(),
	}

	converter.Convert(now, telemetry.GameData{
		Keys: []string{"test", "test2"},
		Data: map[string]float32{
			"test":  1,
			"test2": 123.45,
		},
		RawData: []byte("test,test2\n1,123.45\n"),
	}, 1234)

	fileExists, _ := afero.Exists(fs, "/var/www/simracing-telemetry/test.csv")
	assert.True(t, fileExists)

	theFile, _ := afero.ReadFile(fs, "/var/www/simracing-telemetry/test.csv")
	assert.Equal(t, "test,test2\n1,123.45\n", string(theFile))
}
