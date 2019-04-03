package postgres

import (
	"auth-service/beans"
	"auth-service/config"
	"auth-service/errs"
	"auth-service/models"
	"auth-service/utils"
)

type SessionsRepository struct {
	conn             *Connection
	notAuthorizedErr *errs.Error
}

func NewSessionsRepository(conn *Connection) *SessionsRepository {
	return &SessionsRepository{
		conn:             conn,
		notAuthorizedErr: errs.NewNotAuthorizedError("session not found"),
	}
}

func (r *SessionsRepository) Init() error {
	err := utils.Retry(func(n int) error {
		if n > 0 {
			beans.Logger.Infof("opening storage: retry #%d", n)
		}
		return r.conn.db.Ping()
	}, config.StorageConnNumRet, config.StorageConnRetItv)
	if err != nil {
		return err
	}

	const query = `
        CREATE TABLE IF NOT EXISTS session (
            user_id  TEXT
              CONSTRAINT session_user_id_pk PRIMARY KEY,
            auth_token TEXT NOT NULL,
            refresh_token TEXT NOT NULL,
            expire_time TIMESTAMPTZ NOT NULL
        );
	`

	_, err = r.conn.db.Exec(query)
	return err
}

func (r *SessionsRepository) AddSession(session *models.Session) *errs.Error {
	panic("implement me")
}

func (r *SessionsRepository) UpdateSessionByRefreshToken(session *models.Session) *errs.Error {
	panic("implement me")
}

func (r *SessionsRepository) DeleteSessionByRefreshToken(token string) *errs.Error {
	panic("implement me")
}
