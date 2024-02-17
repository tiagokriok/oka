package storages

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlDB struct {
	*sqlx.DB
}

func NewMysqlDB() (*MysqlDB, error) {
	db, err := sqlx.Connect("mysql", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MysqlDB{db}, err
}

func (db *MysqlDB) Close() {
	db.DB.Close()
}
