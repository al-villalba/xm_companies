package common

import (
	"errors"
	"math"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var database *sql.DB

// GetDatabase : Connect to database
func GetDatabase() *sql.DB {
	if database != nil {
		return database
	}

	e := GetEnv()
	interval := 1000
	attempts := 6
	connCh := make(chan *sql.DB)
	err := errors.New("Connecting")
	var mutex sync.Mutex
	// attempts to connect in increasing intervals 0ms, 1s, 2s, 4s, ...
	for i := 0; i < attempts; i++ {
		go func() {
			sleep := interval * int(math.Floor(math.Pow(float64(2), float64(i-1))))
			time.Sleep(time.Duration(sleep) * time.Millisecond)
			if err == nil {
				// already connected
				return
			}
			mutex.Lock() // ensure one connection attempt at a time
			// logrus.Debugf("Connection attempt %d", i)
			dsn := e.DbUser + ":" + e.DbPass + "@tcp(" + e.DbHost + ")/" + e.DbName
			var db *sql.DB
			db, err = sql.Open("mysql", dsn+"?parseTime=true")
			if err == nil || i == attempts-1 {
				connCh <- db
			}
			mutex.Unlock()
		}()
	}

	database := <-connCh
	if database == nil {
		logrus.WithError(err).Panic("Error connecting to MySQL")
	}

	return database

}
