package db

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/referal_links/internal/storage/postgesql/migration"
)

type Database struct {
	dataSourceName string
	DB             *sqlx.DB
}

func NewDatabase() (*Database, error) {
	var database Database
	config, err := NewDSN()
	if err != nil {
		return nil, err
	}

	err = database.connectDatabase(config)
	if err != nil {
		return nil, err
	}
	err = migration.MigrationsUp(database.DB, database.dataSourceName)
	if err != nil {
		return nil, err
	}

	return &database, err
}

func (db *Database) connectDatabase(config string) error {
	const op = "internal.storage.postgresql.db.ConnectDatabase()"

	db.dataSourceName = config

	var connectError error
	db.DB, connectError = sqlx.Connect("pgx", db.dataSourceName)
	if connectError != nil {
		return fmt.Errorf("func:%s  error:%w", op, connectError)
	}
	return connectError
}

func (db *Database) PingDatabase() error {
	const op = "internal.storage.postgresql.db.PingDatabase()"

	if db.DB == nil {
		return fmt.Errorf("func:%s  error:%s", op, "database isn`t established")
	}

	var pingError = db.DB.Ping()
	if pingError != nil {
		return fmt.Errorf("func:%s  error:%w", op, pingError)
	}
	return nil
}

func (db *Database) CloseDatabase() error {
	const op = "internal.storage.postgresql.db.CloseDatabase()"

	var closingError = db.DB.Close()
	if closingError != nil {
		return fmt.Errorf("func:%s  error:%w", op, closingError)
	}
	return nil
}

func NewDSN() (string, error) {
	const op = "internal.storage.postgresql.db.NewDSN()"

	dsnProtocol := os.Getenv("DB_PROTOCOL")
	if dsnProtocol == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_PROTOCOL isn't set")
	}

	dsnUserName := os.Getenv("DB_USER")
	if dsnUserName == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_USER isn't set")
	}

	dsnPassword := os.Getenv("DB_PASSWORD")
	if dsnPassword == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_PASSWORD isn't set")
	}

	dsnHost := os.Getenv("DB_HOST")
	if dsnHost == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_HOST isn't set")
	}

	dsnPort := os.Getenv("DB_PORT")
	if dsnPort == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_PORT isn't set")
	}

	dsnDBName := os.Getenv("DB_NAME")
	if dsnDBName == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_NAME isn't set")
	}

	dsnOptions := os.Getenv("DB_OPTIONS")
	if dsnOptions == "" {
		return "", fmt.Errorf("func:%s  error:%s", op, "DB_OPTIONS isn't set")
	}

	dsn := dsnProtocol + "://" + dsnUserName + ":" + dsnPassword + "@" +
		dsnHost + ":" + dsnPort + "/" + dsnDBName + "?" + dsnOptions

	return dsn, nil
}
