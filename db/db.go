package db

import (
	"database/sql"
	"log"
	"os"
    "time"

	_ "github.com/mattn/go-sqlite3"
)

var SECURITY_KEY = "xxx"

const (
    DB_FILE              = "./nm.db"
    DB_MAX_OPEN_CONNS    = 60
    DB_MAX_IDLE_CONNS    = 10
    DB_CONN_MAX_LIFETIME = 30 * time.Minute
)

var db *sql.DB

func checkError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func Init(user User) {
    log.Println("INIT DB...")
    var err error
    if _, err := os.Stat(DB_FILE); !os.IsNotExist(err) {
        log.Printf("remove %s\n", DB_FILE)
        os.Remove(DB_FILE)
    }
    db, err = sql.Open("sqlite3", DB_FILE)
    checkError(err)
    defer db.Close()

    log.Println("create tables:")
    for idx, table := range tables {
        log.Printf("\t%d. %s;\n", idx + 1, table.name)
        _, err = db.Exec(table.query)
        checkError(err)
    }

    id, err := AddUser(nil, user)
    checkError(err)
    log.Printf("insert an user[id=%d]\n", id)
}

func Connect() {
    var err error
    db, err = sql.Open("sqlite3", DB_FILE)
    checkError(err)
    db.SetMaxOpenConns(DB_MAX_OPEN_CONNS)
    db.SetMaxIdleConns(DB_MAX_IDLE_CONNS)
    db.SetConnMaxLifetime(DB_CONN_MAX_LIFETIME)
}
