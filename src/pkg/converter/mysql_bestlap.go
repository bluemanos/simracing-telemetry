package converter

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
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

func NewMysqlBestLapConverter(game enums.Game, adapterConfiguration []string) (*MysqlBestLapConverter, error) {
	if len(adapterConfiguration) != 6 {
		return nil, ErrInvalidMySQLAdapterConfiguration
	}
	lastValueCache = make(map[int]hashCache)

	return &MysqlBestLapConverter{
		ConverterData: ConverterData{GameName: game},
		User:          adapterConfiguration[1],
		Password:      adapterConfiguration[2],
		Host:          adapterConfiguration[3],
		Port:          adapterConfiguration[4],
		Database:      adapterConfiguration[5],
		TableName:     "tmd_forzamotorsport2023_bestlaps",
	}, nil
}

// Convert converts the data to the MySQL database
func (db *MysqlBestLapConverter) Convert(_ time.Time, data telemetry.GameData, port int) {
	isBestLap, cacheHash := db.bestLapExists(
		port,
		data.Data["BestLap"],
		data.Data["TrackOrdinal"],
		data.Data["CarOrdinal"],
	)
	if data.Data["BestLap"] == 0 || isBestLap {
		return
	}

	if db.connector == nil {
		fmt.Println("Reconnecting to MySQL BL...")
		var err error
		db.connector, err = sql.Open(
			"mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Password, db.Host, db.Port, db.Database),
		)
		if err != nil {
			log.Println(err)
			return
		}
		db.connector.SetConnMaxLifetime(time.Minute * 5)
		db.connector.SetMaxOpenConns(10)
		db.connector.SetMaxIdleConns(10)
		fmt.Println("Reconnecting to MySQL BL... Connected")
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
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number == 1062 {
			// unique key. Skipping the insert and update the cache
			lastValueCache[port] = cacheHash
			return
		}
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func (db *MysqlBestLapConverter) getHashCacheString(bestLap, trackOrdinal, carOrdinal float32) hashCache {
	return hashCache(fmt.Sprintf("%f-%f-%f", bestLap, trackOrdinal, carOrdinal))
}

func (db *MysqlBestLapConverter) bestLapExists(port int, bestLap, trackOrdinal, carOrdinal float32) (bool, hashCache) {
	hash := db.getHashCacheString(bestLap, trackOrdinal, carOrdinal)
	return lastValueCache[port] == hash, hash
}
