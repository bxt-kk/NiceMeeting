package db

import (
	"database/sql"
	"log"
	"os"
    "time"

	_ "github.com/mattn/go-sqlite3"
)


var (
    SECURITY_KEY         = "xxx"
    DB_FILE              = "./nm.db"
    DB_MAX_OPEN_CONNS    = 60
    DB_MAX_IDLE_CONNS    = 10
    DB_CONN_MAX_LIFETIME = 30 * time.Minute
)

var db *sql.DB

func SetDBFile(filename string) {
    DB_FILE = filename
}

func Init(user User) error {
    log.Println("INIT DB...")
    _, err := os.Stat(DB_FILE);
    if !os.IsNotExist(err) {
        log.Printf("remove %s\n", DB_FILE)
        os.Remove(DB_FILE)
    }
    db, err = sql.Open("sqlite3", DB_FILE)
    if err != nil {
        return err
    }
    defer db.Close()

    log.Println("create tables:")
    for idx, table := range tables {
        log.Printf("\t%d. %s;\n", idx + 1, table.name)
        _, err = db.Exec(table.query)
        if err != nil {
            return err
        }
    }

    id, err := AddUser(nil, user)
    if err != nil {
        return err
    }
    log.Printf("insert an user[id=%d]\n", id)
    return nil
}

func Connect() error {
     _, err := os.Stat(DB_FILE)
     if os.IsNotExist(err) {
         return err
     }
    db, err = sql.Open("sqlite3", DB_FILE)
    if err != nil {
        return err
    }
    db.SetMaxOpenConns(DB_MAX_OPEN_CONNS)
    db.SetMaxIdleConns(DB_MAX_IDLE_CONNS)
    db.SetConnMaxLifetime(DB_CONN_MAX_LIFETIME)
    return nil
}
