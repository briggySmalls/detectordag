package internal

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
)

type Account struct {
    gorm.Model
    Username string
    Password string
    Emails   []Email
    Devices  []Device
}

type Device struct {
    gorm.Model
    ID byte `gorm:"type:BINARY(16)"`
}

type Email struct {
    gorm.Model
    Email string
}

type Database interface {
    Connect(host, port, database, username, password string) error
    Db() *gorm.DB
    Close() error
}

type database struct {
    db *gorm.DB
}

func NewDatabase() {
    return &database{}
}

func Connect(host, port, database, username, password string) error {
    connectionString := fmt.Sprintf(
        "%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
        username,
        password,
        database,
        host,
        port)
    db, err := gorm.Open("mysql", connectionString)
    if err != nil {
        return shared.WrapError(err, "Failed to connect to database")
    }
    return nil
}

func (d *database) Close() error {
    return d.db.Close()
}
