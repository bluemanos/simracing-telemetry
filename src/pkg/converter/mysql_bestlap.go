package converter

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlBestLapConverter struct {
	ConverterData
	User, Password, Host, Port, Database, TableName string
	connector                                       *sql.DB
}

type dbData struct {
	Keys   []string
	Values []string
}

type hashCache string

var lastValueCache map[int]hashCache

// Convert converts the data to the MySQL database
func (db *MysqlBestLapConverter) Convert(_ time.Time, data telemetry.GameData, port int) {
	if db.connector == nil {
		fmt.Println("Reconnecting to MySQL BL...")
		var err error
		db.connector, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Password, db.Host, db.Port, db.Database))
		if err != nil {
			log.Println(err)
			return
		}
	}

	if data.Data["BestLap"] != 0 && db.bestLapExists(port, data.Data["BestLap"], data.Data["TrackOrdinal"], data.Data["CarOrdinal"]) {
		return
	}

	myData := dbData{
		Keys: []string{
			"CarOrdinal", "CarClass", "CarPerformanceIndex", "DrivetrainType", "NumCylinders",
			"Fuel", "BestLap", "LapNumber", "RacePosition", "TrackOrdinal", "user_id",
		},
		Values: []string{
			fmt.Sprintf("%f", data.Data["CarOrdinal"]),
			fmt.Sprintf("%f", data.Data["CarClass"]),
			fmt.Sprintf("%f", data.Data["CarPerformanceIndex"]),
			fmt.Sprintf("%f", data.Data["DrivetrainType"]),
			fmt.Sprintf("%f", data.Data["NumCylinders"]),
			fmt.Sprintf("%f", data.Data["Fuel"]),
			fmt.Sprintf("%f", data.Data["BestLap"]),
			fmt.Sprintf("%f", data.Data["LapNumber"]),
			fmt.Sprintf("%f", data.Data["RacePosition"]),
			fmt.Sprintf("%f", data.Data["TrackOrdinal"]),
			"1",
		},
	}

	values := make([]interface{}, len(myData.Values))
	for i := range myData.Values {
		values[i] = myData.Values[i]
	}

	queryInsertBuilder := sq.Insert(db.TableName).Columns(myData.Keys...).Values(values...)
	query, args, err := queryInsertBuilder.ToSql()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.connector.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return
	}
}

func (db *MysqlBestLapConverter) getHashCacheString(bestLap, trackOrdinal, carOrdinal float32) hashCache {
	return hashCache(fmt.Sprintf("%f-%f-%f", bestLap, trackOrdinal, carOrdinal))
}

func (db *MysqlBestLapConverter) bestLapExists(port int, bestLap, trackOrdinal, carOrdinal float32) bool {
	if lastValueCache == nil {
		lastValueCache = make(map[int]hashCache)
	}

	hash := db.getHashCacheString(bestLap, trackOrdinal, carOrdinal)
	if lastValueCache[port] == hash {
		return true
	}

	lastValueCache[port] = hash
	return false
}
