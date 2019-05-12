package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/beans"
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

	id, err := uuid.NewRandom()
	if err != nil {
		return errs.New(err)
	}

	session.SessionID = id.String()
	_, err = r.conn.db.Exec(query,
		&session.SessionID, &session.UserID, &session.RefreshToken,
	)
	if err != nil {
		return errs.New(err)
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
		return errs.New(err)
	}
	defer r.closeRowsWithCheck(rows)

	if !rows.Next() {
		return r.notAuthorizedErr
	}
	if err := rows.Scan(&session.UserID); err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *SessionsRepository) DeleteSessionByRefreshToken(token string) *errs.Error {
	const query = `
        DELETE FROM session WHERE refresh_token = $1;
    `

	res, err := r.conn.db.Exec(query, &token)
	if err != nil {
		return errs.New(err)
	}
	if n, _ := res.RowsAffected(); n != 1 {
		beans.Logger().Error(fmt.Sprintf("%d sessions deleted", n))
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
		return errs.New(err)
	}

	return nil
}

func (r *SessionsRepository) closeRowsWithCheck(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		beans.Logger().Error(err)
	}
}
