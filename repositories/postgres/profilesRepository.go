package postgres

import (
	"errors"
	"github.com/google/uuid"
	"strings"

	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

type ProfilesRepository struct {
	conn *Connection

	emailDupErr *errs.Error
	notFoundErr *errs.Error
}

func NewProfilesRepository(conn *Connection) *ProfilesRepository {
	return &ProfilesRepository{
		conn: conn,

		emailDupErr: errs.NewConflictError("email duplicate"),
		notFoundErr: errs.NewNotFoundError("profile not found"),
	}
}

func (r *ProfilesRepository) AddProfile(p *models.Profile) *errs.Error {
	const query = `
        INSERT INTO profile(user_id,email,password,is_verified) VALUES($1,$2,$3,false);
    `

	id, err := uuid.NewRandom()
	if err != nil {
		return errs.New(err)
	}
	p.UserID = id.String()

	_, err = r.conn.db.Exec(query,
		p.UserID, p.Credentials.Email, p.Credentials.Password,
	)
	if err != nil {
		if strings.Contains(err.Error(), "profile_email_unique") {
			return r.emailDupErr
		}
		return errs.New(err)
	}

	return nil
}

func (r *ProfilesRepository) SetProfileVerified(p *models.ProfileInfo) *errs.Error {
	const query = `
		UPDATE profile
			SET is_verified = TRUE
			WHERE user_id = $1;
    `

	res, err := r.conn.db.Exec(query, &p.UserID)
	if err != nil {
		return errs.New(err)
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return errs.New(errors.New("no profiles verified"))
	}

	return nil
}

func (r *ProfilesRepository) FindUserIdByCredentials(p *models.Profile) (e *errs.Error) {
	const query = `
        SELECT p.user_id FROM profile p
        WHERE p.email = $1 AND p.password = $2;
    `

	p.IsVerified = true
	rows, err := r.conn.db.Query(query,
		&p.Credentials.Email, &p.Credentials.Password,
	)
	if err != nil {
		return errs.New(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			e = errs.New(err)
		}
	}()

	if !rows.Next() {
		return r.notFoundErr
	}

	if err := rows.Scan(&p.UserID); err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *ProfilesRepository) UpdateEmail(userId, email string) *errs.Error {
	const query = `
		UPDATE profile SET
			email = $2, is_verified = FALSE
		WHERE user_id = $1;
	`

	res, err := r.conn.db.Exec(query, &userId, &email)
	if err != nil {
		return errs.New(err)
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return r.notFoundErr
	}

	return nil
}

func (r *ProfilesRepository) UpdatePassword(userId, password string) *errs.Error {
	panic("implement me") //TODO
}

func (r *ProfilesRepository) DeleteProfileById(p *models.Profile) *errs.Error {
	panic("implement me") //TODO
}
