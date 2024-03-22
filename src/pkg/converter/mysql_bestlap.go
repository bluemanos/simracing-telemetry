package converter

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

var (
	semInsert = semaphore.NewWeighted(1)
	semSelect = semaphore.NewWeighted(1)
)

type MysqlBestLapConverter struct {
	ConverterData
	User, Password, Host, Port, Database, TableName string
	connector                                       *sqlx.DB
}

type dbData struct {
	Keys   []string
	Values []string
}

type BestLapEntity struct {
	ID                  int64     `db:"id"`
	UserID              int64     `db:"user_id"`
	CarOrdinal          int       `db:"CarOrdinal"`
	TrackOrdinal        int       `db:"TrackOrdinal"`
	BestLap             float32   `db:"BestLap"`
	Fuel                float32   `db:"Fuel"`
	CarClass            int       `db:"CarClass"`
	DrivetrainType      int       `db:"DrivetrainType"`
	CarPerformanceIndex int       `db:"CarPerformanceIndex"`
	NumCylinders        int       `db:"NumCylinders"`
	LapNumber           int       `db:"LapNumber"`
	RacePosition        int       `db:"RacePosition"`
	CreatedAt           time.Time `db:"created_at"`
}

type hashCache string

var lastValueCache map[int]map[hashCache]*float32

func NewMysqlBestLapConverter(game enums.Game, adapterConfiguration []string) (*MysqlBestLapConverter, error) {
	if len(adapterConfiguration) != 6 {
		return nil, ErrInvalidMySQLAdapterConfiguration
	}
	lastValueCache = make(map[int]map[hashCache]*float32)

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

func (db *MysqlBestLapConverter) ChannelInit(now time.Time, channel chan telemetry.GameData, port int) {
	fmt.Println("MysqlBestLapConverter ChannelInit")
	for {
		select {
		case data := <-channel:
			db.Convert(now, data, port)
		}
	}
}

// Convert converts the data to the MySQL database
func (db *MysqlBestLapConverter) Convert(_ time.Time, data telemetry.GameData, port int) {
	if db.connector == nil {
		fmt.Println("Reconnecting to MySQL BL...")
		var err error
		db.connector, err = sqlx.Open(
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

	isBestLap, cacheHash := db.bestLapExists(
		port,
		data.Data["BestLap"],
		data.Data["TrackOrdinal"],
		data.Data["CarOrdinal"],
	)
	if data.Data["BestLap"] == 0 || !isBestLap {
		return
	}
	
	if !semInsert.TryAcquire(1) {
		return
	}
	defer semInsert.Release(1)

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
	telemetry.DisplayLog("vvv", query)
	telemetry.DisplayLog("vvv", args)

	_, err = db.connector.Exec(query, args...)
	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) {
		if mysqlError.Number == 1062 {
			// unique key. Skipping the insert and update the cache
			currentBestLap := data.Data["BestLap"]
			lastValueCache[port][cacheHash] = &currentBestLap
			return
		}
	}
	if err != nil {
		log.Println(err)
		return
	}

	currentBestLap := data.Data["BestLap"]
	lastValueCache[port][cacheHash] = &currentBestLap
}

func (db *MysqlBestLapConverter) getHashCacheString(bestLap, trackOrdinal, carOrdinal float32) hashCache {
	return hashCache(fmt.Sprintf("%f-%f-%f", bestLap, trackOrdinal, carOrdinal))
}

func (db *MysqlBestLapConverter) bestLapExists(port int, bestLap, trackOrdinal, carOrdinal float32) (bool, hashCache) {
	hash := db.getHashCacheString(bestLap, trackOrdinal, carOrdinal)

	if bestLap == 0 {
		return false, hash
	}

	if lastValueCache[port] == nil {
		lastValueCache[port] = make(map[hashCache]*float32)
	}

	if lastValueCache[port][hash] == nil {
		if !semSelect.TryAcquire(1) {
			return false, hash
		}
		defer semSelect.Release(1)

		queryInsertBuilder := sq.Select([]string{"id", "BestLap"}...).
			From("tmd_forzamotorsport2023_bestlaps").
			Where(sq.Eq{
				"TrackOrdinal": trackOrdinal,
				"CarOrdinal":   carOrdinal,
				"user_id":      "1",
			}).
			OrderBy("BestLap ASC").
			Limit(1)

		query, args, err := queryInsertBuilder.ToSql()
		if err != nil {
			log.Println(err)
			return false, hash
		}
		telemetry.DisplayLog("vvv", query)
		telemetry.DisplayLog("vvv", args)

		bestLapDb := BestLapEntity{}
		err = db.connector.Get(&bestLapDb, query, args...)
		if errors.Is(err, sql.ErrNoRows) {
			noBestLap := float32(0)
			lastValueCache[port][hash] = &noBestLap
		}
		if err != nil {
			return true, hash
		}

		lastValueCache[port][hash] = &bestLapDb.BestLap
	}

	return *lastValueCache[port][hash] > bestLap, hash
}
