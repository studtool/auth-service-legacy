package postgres

import (
	"auth-service/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Connection struct {
	connStr string
	db      *sql.DB
}

func NewConnection() *Connection {
	return &Connection{
		connStr: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.StorageUser, config.StoragePassword,
			config.StorageHost, config.StoragePort,
			config.StorageDB, config.StorageSSL,
		),
	}
}

func (conn *Connection) Open() (err error) {
	conn.db, err = sql.Open("postgres", conn.connStr)
	return err
}

func (conn *Connection) Close() (err error) {
	return conn.db.Close()
}
