package internal

import (
    "database/sql"
    "fmt"
    "github.com/briggysmalls/detectordag/shared"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type DbParams struct {
    Host     string
    Port     int32
    Database string
    Username string
    Password string
}

type Database interface {
    Connect(params DbParams) error
    DB() *sql.DB
    Close() error
}

type database struct {
    db *sql.DB
}

func NewDatabase() Database {
    return &database{}
}

func (d *database) Connect(params DbParams) error {
    connectionString := fmt.Sprintf(
        "%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
        params.Username,
        params.Password,
        params.Host,
        params.Port,
        params.Database)
    log.Printf("Connecting to database: %s", connectionString)
    db, err := sql.Open("mysql", connectionString)
    if err != nil {
        return shared.WrapError(err, "Failed to connect to database")
    }
    d.db = db
    return nil
}

func (d *database) DB() *sql.DB {
    return d.db
}

func (d *database) Close() error {
    return d.db.Close()
}
