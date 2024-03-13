package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

type Config interface {
	GetDB() *sql.DB
}

type Database struct {
	db  *sql.DB
}

var onceLoadDb sync.Once

func (dbCon *Database) GetDB() *sql.DB {
	onceLoadDb.Do(func() {
		db, err := sql.Open("postgres", "user=postgres host=localhost password=root123 dbname=service-user sslmode=disable")
		if err != nil {
			log.Fatal("Cannot start app, Error when connect to DB ", err.Error())
		}

		err = db.Ping()
		if err != nil {
		panic(err)
		}
		fmt.Println("Successfully connected to PostgreSQL database!")

		dbCon.db = db
	})
	return dbCon.db
}

func NewInfraManager() Config {
	infra := Database{}
	return &infra
}
