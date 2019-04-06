package postgres

import (
	"auth-service/models"
	"github.com/hashicorp/go-uuid"
	"github.com/studtool/common/errs"
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

func (r *ProfilesRepository) FindProfileByCredentials(p *models.Profile) *errs.Error {
	const query = `
        SELECT p.user_id FROM profile p
        WHERE p.email = $1 AND p.password = $2;
    `

	panic("implement me") //TODO
}

func (r *ProfilesRepository) UpdateCredentials(c *models.Credentials) *errs.Error {
	panic("implement me") //TODO
}

func (r *ProfilesRepository) DeleteProfileById(p *models.Profile) *errs.Error {
	panic("implement me") //TODO
}
