package connection

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config interface {
	GetDB() *sql.DB
}

type Database struct {
	db *sql.DB
}

func init() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

func (dbCon *Database) GetDB() *sql.DB {
	var db *sql.DB
	var err error

	for i := 1; i <= 20; i++ {
		connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_PORT"), 
		os.Getenv("DB_NAME"))

		db, err = sql.Open("postgres", connStr)
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
