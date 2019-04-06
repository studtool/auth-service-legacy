package postgres

import (
	"github.com/hashicorp/go-uuid"

	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
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
	const query = `
        INSERT INTO session(session_id, user_id, refresh_token) VALUES ($1,$2,$3);
    `

	sessionId, _ := uuid.GenerateUUID()

	_, err := r.conn.db.Exec(query,
		&sessionId, &session.UserId, &session.RefreshToken,
	)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}

func (r *SessionsRepository) FindUserIdByRefreshToken(session *models.Session) (e *errs.Error) {
	const query = `
        SELECT s.user_id FROM session s
        WHERE s.refresh_token = $1;
    `

	rows, err := r.conn.db.Query(query,
		&session.RefreshToken,
	)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	defer func() {
		if err := rows.Close(); err != nil {
			e = errs.NewInternalError(err.Error())
		}
	}()

	if !rows.Next() {
		return r.notAuthorizedErr
	}

	if err := rows.Scan(&session.UserId); err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}

func (r *SessionsRepository) DeleteSessionByRefreshToken(token string) *errs.Error {
	const query = `
        DELETE FROM session WHERE refresh_token = $1;
    `

	if _, err := r.conn.db.Exec(query, &token); err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}

func (r *SessionsRepository) DeleteAllSessionsByRefreshToken(token string) *errs.Error {
	const query = `
        DELETE FROM session WHERE user_id = (
            SELECT s.user_id FROM session s
            WHERE s.refresh_token = $1
        );
    `

	_, err := r.conn.db.Exec(query, &token)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}
