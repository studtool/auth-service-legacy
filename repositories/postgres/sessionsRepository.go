package postgres

import (
	"auth-service/errs"
	"auth-service/models"
	"auth-service/types"
	"context"
	"database/sql"
	"fmt"
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

func (r *SessionsRepository) AddSession(session *models.Session) (res *errs.Error) {
	tx, err := r.conn.db.BeginTx(context.TODO(), &sql.TxOptions{})
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	defer func() {
		if err := tx.Commit(); err != nil {
			res = errs.NewInternalError(err.Error())
		}
	}()

	const selectQuery = `
        SELECT EXISTS (SELECT * FROM session WHERE user_id=$1);
    `

	rows, err := tx.Query(selectQuery, &session.UserId)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	defer func() {
		if err := rows.Close(); err != nil {
			res = errs.NewInternalError(err.Error())
		}
	}()

	if !rows.Next() {
		errs.NewInternalError(fmt.Sprintf(`!rows.Next() in "%s"`, selectQuery))
	}

	var exists bool
	if err := rows.Scan(&exists); err != nil {
		return errs.NewInternalError(err.Error())
	}

	if exists {
		const query = `
            SELECT expire_time FROM session WHERE user_id=$1;
        `

		rows, err := tx.Query(query, &session.UserId)
		if err != nil {
			return errs.NewInternalError(err.Error())
		}
		defer func() {
			if err := rows.Close(); err != nil {
				res = errs.NewInternalError(err.Error())
			}
		}()

		if !rows.Next() {
			errs.NewInternalError(fmt.Sprintf(`!rows.Next() in "%s"`, query))
		}

		var expTime time.Time
		if err := rows.Scan(&expTime); err != nil {
			return errs.NewInternalError(err.Error())
		}

		session.ExpireTime = types.DateTime(expTime)
	} else {
		const query = `
            INSERT INTO session(user_id,auth_token,refresh_token,expire_time) VALUES($1,$2,$3,$4);
        `

		expTime := time.Time(session.ExpireTime)
		res, err := tx.Exec(query,
			&session.UserId, &session.AuthToken, &session.RefreshToken, &expTime,
		)
		if err != nil {
			return errs.NewInternalError(err.Error())
		}
		if n, _ := res.RowsAffected(); n != 1 {
			return errs.NewInternalError(fmt.Sprintf(`(res.RowsAffected() != 1) in "%s"`, query))
		}
	}

	return nil
}

func (r *SessionsRepository) UpdateSessionByRefreshToken(session *models.Session) *errs.Error {
	panic("implement me")
}

func (r *SessionsRepository) DeleteSessionByRefreshToken(token string) *errs.Error {
	panic("implement me")
}
