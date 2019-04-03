package postgres

import (
	"auth-service/errs"
	"auth-service/models"
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

func (r *SessionsRepository) AddSession(session *models.Session) *errs.Error {
	panic("implement me")
}

func (r *SessionsRepository) UpdateSessionByRefreshToken(session *models.Session) *errs.Error {
	panic("implement me")
}

func (r *SessionsRepository) DeleteSessionByRefreshToken(token string) *errs.Error {
	panic("implement me")
}
