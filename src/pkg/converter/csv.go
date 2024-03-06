package converter

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/spf13/afero"
)

var (
	ErrInvalidCsvAdapterConfiguration = errors.New("[CSV] invalid adapter configuration")
	ErrInvalidFilePath                = errors.New("[CSV] invalid file path")
	ErrInvalidRetention               = errors.New("[CSV] invalid retention type")
)

type CsvConverter struct {
	ConverterData
	Fs          afero.Fs
	FilePath    string
	Retention   enums.RetentionType
	fileHandler afero.File
}

func NewCsvConverter(game enums.Game, adapterConfiguration []string, fs afero.Fs) (*CsvConverter, error) {
	if len(adapterConfiguration) != 3 {
		log.Printf("[%s] Wrong CSV adapter configuration", game)
		return nil, ErrInvalidCsvAdapterConfiguration
	}

	return &CsvConverter{
		ConverterData: ConverterData{GameName: game},
		Fs:            fs,
		FilePath:      adapterConfiguration[1],
		Retention:     enums.RetentionType(adapterConfiguration[2]),
	}, nil
}

// Convert the data to CSV format and writes it to the file
func (csv *CsvConverter) Convert(now time.Time, data telemetry.GameData, _ int) {
	afs := &afero.Afero{Fs: csv.Fs}
	filePath, err := csv.CorrectFilePath(now)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fileExists, err := afs.Exists(filePath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if !fileExists {
		csvHeader := ""
		for _, key := range data.Keys {
			csvHeader += fmt.Sprintf(",%s", key)
		}
		csvHeader = csvHeader + "\n"
		err = afs.WriteFile(filePath, []byte(csvHeader)[1:], 0o644)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	if csv.fileHandler == nil || csv.fileHandler.Name() != filePath {
		csv.fileHandler, err = afs.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer csv.fileHandler.Close()
	}

	csvLine := ""
	for _, key := range data.Keys {
		csvLine += fmt.Sprintf(",%v", data.Data[key])
	}
	csvLine += "\n"
	fmt.Fprintf(csv.fileHandler, csvLine[1:])
}

// CorrectFilePath returns the correct file path based on the retention type
func (csv *CsvConverter) CorrectFilePath(now time.Time) (string, error) {
	afs := &afero.Afero{Fs: csv.Fs}

	switch csv.Retention {
	case enums.RetentionTypes.Daily():
		return csv.dailyRetention(now, afs)
	case enums.RetentionTypes.None():
		return csv.noRetention(now, afs)
	}

	return "", ErrInvalidRetention
}

// dailyRetention validate and returns the file path for daily retention
func (csv *CsvConverter) dailyRetention(now time.Time, afs *afero.Afero) (string, error) {
	isDir, err := afs.IsDir(csv.FilePath)
	if err != nil || !isDir {
		return "", ErrInvalidFilePath
	}

	defaultFileName := fmt.Sprintf("%s-daily-%s.csv", csv.GameName, now.Format("2006-01-02"))

	slashAtTheEnd := csv.FilePath[len(csv.FilePath)-1:]
	if slashAtTheEnd != "/" {
		csv.FilePath += "/"
	}

	return csv.FilePath + defaultFileName, nil
}

// noRetention validate and returns the file path for no retention type
func (csv *CsvConverter) noRetention(_ time.Time, afs *afero.Afero) (string, error) {
	defaultFileName := fmt.Sprintf("%s.csv", csv.GameName)

	dir, file := filepath.Split(csv.FilePath)

	if file == "" {
		isDir, err := afs.IsDir(dir)
		if err != nil || !isDir {
			return "", ErrInvalidFilePath
		}

		return csv.FilePath + defaultFileName, nil
	}

	fileExt := filepath.Ext(csv.FilePath)
	isDir, err := afs.IsDir(csv.FilePath)
	if fileExt != ".csv" && (err != nil || !isDir) {
		return "", ErrInvalidFilePath
	}

	if fileExt == ".csv" {
		return csv.FilePath, nil
	}

	slashAtTheEnd := csv.FilePath[len(csv.FilePath)-1:]
	if slashAtTheEnd != "/" {
		csv.FilePath += "/"
	}

	return csv.FilePath + defaultFileName, nil
}
