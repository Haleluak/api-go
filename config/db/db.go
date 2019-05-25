package db

import (
	"database/sql"
	"errors"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type DataBaseMysql struct {
	HostMysql string
	UserMysql string
	PassMysql string
	DBMysql   string
}

func (db *DataBaseMysql) ConnectDB() (*sql.DB, error) {
	DBStringConnect := db.UserMysql + ":" + db.PassMysql + "@tcp(" + db.HostMysql + ")/" + db.DBMysql
	dbClient, err := sql.Open("mysql", DBStringConnect)
	if err != nil {
		log.Println(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return nil, errors.New("Error Connect Database")
	}
	return dbClient, nil
}

