package converter

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bluemanos/simracing-telemetry/src/pkg/enums"
	"github.com/bluemanos/simracing-telemetry/src/telemetry"
	_ "github.com/go-sql-driver/mysql"
)

var ErrInvalidMySQLAdapterConfiguration = errors.New("[MySQL] invalid adapter configuration")

type MySQLConverter struct {
	ConverterData
	User, Password, Host, Port, Database, TableName string
	connector                                       *sql.DB
}

func NewMySQLConverter(game enums.Game, adapterConfiguration []string) (*MySQLConverter, error) {
	if len(adapterConfiguration) != 6 {
		return nil, ErrInvalidMySQLAdapterConfiguration
	}

	return &MySQLConverter{
		ConverterData: ConverterData{GameName: game},
		User:          adapterConfiguration[1],
		Password:      adapterConfiguration[2],
		Host:          adapterConfiguration[3],
		Port:          adapterConfiguration[4],
		Database:      adapterConfiguration[5],
		TableName:     gameEnvKeys[game].DatabaseTable,
	}, nil
}

// Convert converts the data to the MySQL database
func (db *MySQLConverter) Convert(_ time.Time, data telemetry.GameData, _ int) {
	if data.Data["IsRaceOn"] == 0 {
		return
	}

	if db.connector == nil {
		fmt.Println("Reconnecting to MySQL...")
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
	}

	values := make([]interface{}, len(data.Keys))
	for i, key := range data.Keys {
		values[i] = data.Data[key]
	}

	queryInsertBuilder := sq.Insert(db.TableName).Columns(data.Keys...).Values(values...)
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
