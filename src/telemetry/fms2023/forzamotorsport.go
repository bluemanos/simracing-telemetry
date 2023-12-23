package fms2023

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/bluemanos/simracing-telemetry/src/pkg/server"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	_ "github.com/go-sql-driver/mysql"
)

const dataFormatFile = "forzamotorsport"

type ForzaMotorsportHandler struct {
	telemetry.TelemetryHandler
	DebugMode string
}

func (fm *ForzaMotorsportHandler) InitAndRun(port int) error {
	udpServer := server.UDPServer{
		Addr: "0.0.0.0:" + strconv.Itoa(port),
	}

	fm.Telemetries, fm.Keys = fm.initTelemetry()

	log.Printf("Forza data out server listening on %s:%d, waiting for Forza data...\n", telemetry.GetOutboundIP(), port)

	err := udpServer.Run(fm.processBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (fm *ForzaMotorsportHandler) initTelemetry() (map[string]telemetry.TelemetryData, []string) {
	lines, err := telemetry.ReadLines("fms2023/" + dataFormatFile)
	if err != nil {
		log.Fatalf("Error reading format file: %s", err)
	}

	telemetryArray := make(map[string]telemetry.TelemetryData, len(lines))
	telemetryKeys := make([]string, len(lines))
	startOffset := 0
	endOffset := 0
	dataLength := 0

	for i, line := range lines {
		dataFormat := strings.Split(line, " ")
		dataType := dataFormat[0]
		dataName := dataFormat[1]

		switch dataType {
		case "S32", "U32", "F32":
			dataLength = 4
		case "U16":
			dataLength = 2
		case "U8", "S8":
			dataLength = 1
		default:
			log.Fatalf("ForzaMotorsportHandler Error: Unknown data type: %s\n", dataType)
		}
		endOffset = endOffset + dataLength
		startOffset = endOffset - dataLength

		telemItem := telemetry.TelemetryData{
			Position:    i,
			Name:        dataName,
			DataType:    dataType,
			StartOffset: startOffset,
			EndOffset:   endOffset,
		}
		telemetryArray[dataName] = telemItem
		telemetryKeys[i] = dataName
	}

	return telemetryArray, telemetryKeys
}

func (fm *ForzaMotorsportHandler) processBuffer(buffer []byte) {
	for i, telemetryObj := range fm.Telemetries {
		data := buffer[telemetryObj.StartOffset:telemetryObj.EndOffset]

		switch telemetryObj.DataType {
		case "F32":
			telemetryObj.Data = math.Float32frombits(binary.LittleEndian.Uint32(data))
		case "U8":
			telemetryObj.Data = float32(data[0])
		case "S8":
			telemetryObj.Data = float32(int8(data[0]))
		case "U16":
			telemetryObj.Data = float32(binary.LittleEndian.Uint16(data))
		default:
			telemetryObj.Data = float32(binary.LittleEndian.Uint32(data))
		}

		fm.Telemetries[i] = telemetryObj
	}

	telemetry.DisplayLog("vvv", fmt.Sprintf(
		"IsRace: %.0f \t RPM: %.0f \t Gear: %.0f \t BHP: %.0f \t Speed: %.0f \t Total slip: %.0f",
		fm.Telemetries["IsRaceOn"].Data,
		fm.Telemetries["CurrentEngineRpm"].Data,
		fm.Telemetries["Gear"].Data,
		math.Max(0.0, float64(fm.Telemetries["Power"].Data/745.699872)),
		fm.Telemetries["Speed"].Data*3.6, // 3.6 for kph, 2.237 for mph
		fm.Telemetries["TireCombinedSlipRearLeft"].Data+fm.Telemetries["TireCombinedSlipRearRight"].Data,
	))

	//db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/app")
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer db.Close()
	//
	//query := fmt.Sprintf("INSERT INTO tmd_forzamotorsport2023 (IsRaceOn, TimestampMS, EngineMaxRpm, EngineIdleRpm, CurrentEngineRpm) VALUES(%.4f, %.4f, %.4f, %.4f, %.4f)", fm.telemetries["IsRaceOn"].data, fm.telemetries["TimestampMS"].data, fm.telemetries["EngineMaxRpm"].data, fm.telemetries["EngineIdleRpm"].data, fm.telemetries["CurrentEngineRpm"].data)
	//insert, err := db.Query(query)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer insert.Close()
}
