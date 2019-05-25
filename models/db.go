package models

import (
	"fmt"
	"github.com/Haleluak/kb-backend/config/db"
	"github.com/joho/godotenv"
	"os"
)

var DatabaseMysql *db.DataBaseMysql

func InitDB() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	DatabaseMysql = &db.DataBaseMysql{
	HostMysql: os.Getenv("db_host"),
	UserMysql: os.Getenv("db_user"),
	PassMysql: os.Getenv("db_pass"),
	DBMysql:   os.Getenv("db_name"),
	}
}
