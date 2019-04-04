package postgres

import (
	"auth-service/models"
	"github.com/studtool/common/errs"
	"time"
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

func (r *SessionsRepository) AddSession(credentials *models.Credentials, session *models.Session) (res *errs.Error) {
	const query = `
        SELECT add_session($1,$2,$3,$4,$5);
    `

	t := time.Time(session.ExpireTime)

	rows, err := r.conn.db.Query(query,
		&credentials.Email, &credentials.Password,
		&session.AuthToken, &session.RefreshToken, &t,
	)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	defer func() {
		if err := rows.Close(); err != nil {
			res = errs.NewInternalError(err.Error())
		}
	}()

	if !rows.Next() {
		return r.notAuthorizedErr
	}

	var userId *string
	if err := rows.Scan(&userId); err != nil {
		return errs.NewInternalError(err.Error())
	}
	if userId == nil {
		return r.notAuthorizedErr
	}

	session.UserId = *userId
	return nil
}

func (r *SessionsRepository) UpdateSessionByRefreshToken(session *models.Session) *errs.Error {
	panic("implement me")
}

func (r *SessionsRepository) DeleteSessionByRefreshToken(token string) *errs.Error {
	panic("implement me")
}
