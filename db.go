package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
)

type Log struct {
	Tz       time.Time `json:"tz,omitempty"`
	Envelope Envelope  `json:"envelope,omitempty"`
}

func EnvelopeToLog(e Envelope) Log {
	return Log{
		Tz:       time.Now(),
		Envelope: e,
	}
}

type Database struct {
	Connection redis.Conn
	Handler    *rejson.Handler
}

func InitializeDatabase(db *Database) error {
	emptyArray := []Envelope{}
	_, err := db.Handler.JSONSet("logs", "$", emptyArray)
	return err
}

func EnsureDatabaseIsInitialized(db *Database) (bool, error) {
	keyExists := true
	resp, err := db.Handler.JSONGet("logs", "$")

	//	@todo: check the API for what actually get returns for "key doesn't exist"
	if err != nil || resp == nil {
		keyExists = false
		err = InitializeDatabase(db)
	}

	return keyExists, err
}

func (db *Database) Connect(connectionString string) error {
	var addr = &connectionString
	rh := rejson.NewReJSONHandler()
	conn, err := redis.Dial("tcp", *addr)
	if err != nil {
		return fmt.Errorf("failed to connect to redis-stack-server @ %s", *addr)
	}
	rh.SetRedigoClient(conn)
	db.Connection = conn
	db.Handler = rh
	return nil
}

func (db *Database) Disconnect() error {
	_, err := db.Connection.Do("FLUSHALL")
	if err == nil {
		err = db.Connection.Close()
	}
	return err
}

func NewDatabaseWithConnection(connectionString string) (*Database, error) {
	db := Database{}
	err := db.Connect(connectionString)
	return &db, err
}

func (db *Database) AllLogs() ([]Log, error) {
	var err error
	var logs []Log
	res, err := db.Handler.JSONGet("logs", "$")
	if err != nil {
		return logs, err
	}
	// var ArrOut []string
	// err = json.Unmarshal(res.([]byte), &ArrOut)
	// if err != nil {
	// 	return logs, fmt.Errorf("failed to JSON Unmarshal")
	// }
	err = json.Unmarshal(res.([]byte), &logs)
	return logs, err
}

func (db *Database) Log(e Envelope) error {
	var err error = nil
	_, err = db.Handler.JSONArrAppend("logs", "$", EnvelopeToLog(e))
	//	@todo: maybe examine actual response
	return err
}
