package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/studtool/auth-service/config"
)

type Connection struct {
	connStr string
	db      *sql.DB
}

func NewConnection() *Connection {
	return &Connection{
		connStr: fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			config.AuthStorageUser.Value(), config.AuthStoragePassword.Value(),
			config.AuthStorageHost.Value(), config.AuthStoragePort.Value(),
			config.AuthStorageName.Value(), config.AuthStorageSSL.Value(),
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
