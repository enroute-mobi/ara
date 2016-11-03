package model

import (
	"database/sql"
	"fmt"

	"github.com/af83/edwig/logger"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

var Database *gorp.DbMap

type DatabaseConfig struct {
	Name     string
	User     string
	Password string
	Port     uint
}

func InitDB(config DatabaseConfig) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.User,
		config.Password,
		config.Name,
	)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		logger.Log.Panicf("Error while connecting to the database:\n%v", err)
	}
	logger.Log.Debugf("Connected to Database %s", config.Name)
	// construct a gorp DbMap
	Database = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
}

func CloseDB() {
	Database.Db.Close()
}
