package postgres

import (
	"github.com/hashicorp/go-uuid"
	"strings"

	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
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

func (r *ProfilesRepository) FindUserIdByCredentials(p *models.Profile) (e *errs.Error) {
	const query = `
        SELECT p.user_id FROM profile p
        WHERE p.email = $1 AND p.password = $2;
    `

	rows, err := r.conn.db.Query(query,
		&p.Credentials.Email, &p.Credentials.Password,
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
		return r.notFoundErr
	}

	if err := rows.Scan(&p.UserId); err != nil {
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
