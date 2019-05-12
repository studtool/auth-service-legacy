package postgres

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/beans"
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

	hash, err := r.getPasswordHash(p.Password)
	if err != nil {
		return errs.New(err)
	}

	_, err = r.conn.db.Exec(query,
		p.UserID, p.Email, hash,
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

func (r *ProfilesRepository) FindUserIdByCredentials(p *models.Profile) *errs.Error {
	const query = `
		SELECT
			p.user_id,
			p.password
		FROM profile p
        WHERE p.email = $1 AND is_verified;
    `

	rows, err := r.conn.db.Query(query, &p.Email)
	if err != nil {
		return errs.New(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			beans.Logger.Error(errs.New(err)) //TODO format
		}
	}()

	if !rows.Next() {
		return r.notFoundErr
	}

	var password string
	if err := rows.Scan(&p.UserID, &password); err != nil {
		return errs.New(err)
	}
	if err := r.checkPassword(p.Password, password); err != nil {
		return r.notFoundErr
	}

	return nil
}

func (r *ProfilesRepository) UpdateEmail(u *models.EmailUpdate) *errs.Error {
	const query = `
		UPDATE profile SET
			email = $2, is_verified = FALSE
		WHERE user_id = $1 AND is_verified;
	`

	res, err := r.conn.db.Exec(query, &u.UserID, &u.Email)
	if err != nil {
		return errs.New(err)
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return r.notFoundErr
	}

	return nil
}

func (r *ProfilesRepository) UpdatePassword(u *models.PasswordUpdate) *errs.Error {
	const query = `
		UPDATE profile SET
			password = $2, is_verified = FALSE
		WHERE user_id = $1 AND is_verified;
	`

	res, err := r.conn.db.Exec(query, &u.UserID, &u.Password)
	if err != nil {
		return errs.New(err)
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return r.notFoundErr
	}

	return nil
}

func (r *ProfilesRepository) DeleteProfileById(userId string) *errs.Error {
	const query = `
		DELETE FROM profile
		WHERE user_id = $1 AND is_verified;
	`

	res, err := r.conn.db.Exec(query, &userId)
	if err != nil {
		return errs.New(err)
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return r.notFoundErr
	}

	return nil
}

func (r *ProfilesRepository) getPasswordHash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //TODO optimize
	if err != nil {
		return consts.EmptyString, errs.New(err)
	}
	return string(h), nil //TODO optimize
}

func (r *ProfilesRepository) checkPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
