package connection

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"
)

type Config interface {
	GetDB() *sql.DB
}

type Database struct {
	db *sql.DB
}

func (dbCon *Database) GetDB() *sql.DB {
	var db *sql.DB
	var err error

	for i := 1; i <= 20; i++ {
		db, err = sql.Open("postgres", "user=postgresapi host=db port=5432 password=root dbname=apigateaway sslmode=disable")
		if err == nil {
			err = db.Ping()
		}

		if err != nil { 
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Op == "dial" {
					time.Sleep(time.Duration(i) * 100 * time.Millisecond)
					continue
				} else {
					log.Fatalf("Cannot start app, Error when connect to DB %v", err.Error())
				}
			} else {
				log.Fatalf("Cannot start app, Error when connect to DB %v", err.Error())
			}
		}

		break
	}

	if err != nil {
		log.Fatalf("Cannot start app, Error when connect to DB %v", err.Error())
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	dbCon.db = db

	return dbCon.db
}

func NewPostgres() Config {
	infra := Database{}
	return &infra
}
