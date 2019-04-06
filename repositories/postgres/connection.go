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
			config.StorageUser.Value(), config.StoragePassword.Value(),
			config.StorageHost.Value(), config.StoragePort.Value(),
			config.StorageDB.Value(), config.StorageSSL.Value(),
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
