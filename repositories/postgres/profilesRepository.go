package postgres

import (
	"auth-service/errs"
	"auth-service/models"
	"github.com/hashicorp/go-uuid"
	"strings"
)

type ProfilesRepository struct {
	conn        *Connection
	notFoundErr *errs.Error
}

func NewProfilesRepository(conn *Connection) *ProfilesRepository {
	return &ProfilesRepository{
		conn:        conn,
		notFoundErr: errs.NewNotFoundError("profile not found"),
	}
}

func (r *ProfilesRepository) AddProfile(p *models.Profile) *errs.Error {
	const query = `
        INSERT INTO profile(user_id,email,password) VALUES($1,$2,$3);
    `

	p.UserId, _ = uuid.GenerateUUID()

	_, err := r.conn.db.Exec(query,
		p.UserId, p.Credentials.Email, p.Credentials.Password,
	)
	if err != nil {
		if strings.Contains(err.Error(), "profile_user_id_pk") {
			return errs.NewInternalError(err.Error())
		}
		if strings.Contains(err.Error(), "profile_email_unique") {
			return errs.NewConflictError("email duplicate")
		}
		return errs.NewInternalError(err.Error())
	}

	return nil
}

func (r *ProfilesRepository) GetProfileIDByCredentials(p *models.Profile) *errs.Error {
	const query = `
        SELECT user_id FROM profile WHERE email=$1 AND password=$2;
    `

	row, err := r.conn.db.Query(query,
		&p.Credentials.Email, &p.Credentials.Password,
	)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	defer func() {
		if err := row.Close(); err != nil {
			panic(err)
		}
	}()

	if !row.Next() {
		return r.notFoundErr
	}
	if err := row.Scan(&p.UserId); err != nil {
		return errs.NewInternalError(err.Error())
	}

	return nil
}

func (r *ProfilesRepository) UpdateCredentials(c *models.Credentials) *errs.Error {
	panic("implement me") //TODO
}

func (r *ProfilesRepository) DeleteProfileById(p *models.Profile) *errs.Error {
	panic("implement me") //TODO
}
