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

type MySqlConverter struct {
	ConverterData
	User, Password, Host, Port, Database, TableName string
	connector                                       *sql.DB
}

// Convert converts the data to the MySQL database
func (db *MySqlConverter) Convert(_ time.Time, data telemetry.GameData, port int) {
	if data.Data["IsRaceOn"] == 0 {
		return
	}

	if db.connector == nil {
		fmt.Println("Reconnecting to MySQL...")
		var err error
		db.connector, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Password, db.Host, db.Port, db.Database))
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
