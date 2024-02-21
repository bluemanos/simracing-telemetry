package converter

import (
	"errors"
	"fmt"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/spf13/afero"
)

var (
	errInvalidFilePath  = errors.New("[CSV] invalid file path")
	errInvalidRetention = errors.New("[CSV] invalid retention type")
)

type CsvConverter struct {
	ConverterData
	Fs          afero.Fs
	FilePath    string
	Retention   enums.RetentionType
	fileHandler afero.File
}

// Convert the data to CSV format and writes it to the file
func (csv CsvConverter) Convert(now time.Time, data telemetry.GameData) {
	afs := &afero.Afero{Fs: csv.Fs}
	filePath, err := csv.correctFilePath(now)
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
		err = afs.WriteFile(filePath, []byte(csvHeader)[1:], 0644)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	if csv.fileHandler == nil || csv.fileHandler.Name() != filePath {
		csv.fileHandler, err = afs.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
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

// correctFilePath returns the correct file path based on the retention type
func (csv CsvConverter) correctFilePath(now time.Time) (string, error) {
	afs := &afero.Afero{Fs: csv.Fs}

	switch csv.Retention {
	case enums.RetentionTypes.Daily():
		return csv.dailyRetention(now, afs)
	case enums.RetentionTypes.None():
		return csv.noRetention(now, afs)
	}

	return "", errInvalidRetention
}

// dailyRetention validate and returns the file path for daily retention
func (csv CsvConverter) dailyRetention(now time.Time, afs *afero.Afero) (string, error) {
	isDir, err := afs.IsDir(csv.FilePath)
	if err != nil || !isDir {
		return "", errInvalidFilePath
	}

	defaultFileName := fmt.Sprintf("%s-daily-%s.csv", csv.GameName, now.Format("2006-01-02"))

	slashAtTheEnd := csv.FilePath[len(csv.FilePath)-1:]
	if slashAtTheEnd != "/" {
		csv.FilePath += "/"
	}

	return csv.FilePath + defaultFileName, nil
}

// noRetention validate and returns the file path for no retention type
func (csv CsvConverter) noRetention(_ time.Time, afs *afero.Afero) (string, error) {
	defaultFileName := fmt.Sprintf("%s.csv", csv.GameName)

	dir, file := filepath.Split(csv.FilePath)

	if file == "" {
		isDir, err := afs.IsDir(dir)
		if err != nil || !isDir {
			return "", errInvalidFilePath
		}

		return csv.FilePath + defaultFileName, nil
	}

	fileExt := filepath.Ext(csv.FilePath)
	isDir, err := afs.IsDir(csv.FilePath)
	if fileExt != ".csv" && (err != nil || !isDir) {
		return "", errInvalidFilePath
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
