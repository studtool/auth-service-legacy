package postgres

import (
	"auth-service/beans"
	"auth-service/config"
	"auth-service/errs"
	"auth-service/models"
	"auth-service/utils"
	"github.com/hashicorp/go-uuid"
	"strings"
)

type ProfilesRepository struct {
	conn *Connection
}

func NewProfilesRepository(conn *Connection) *ProfilesRepository {
	return &ProfilesRepository{
		conn: conn,
	}
}

func (r *ProfilesRepository) Init() error {
	err := utils.Retry(func(n int) error {
		if n > 0 {
			beans.Logger.Infof("opening storage: retry #%d", n)
		}
		return r.conn.db.Ping()
	}, config.StorageConnNumRet, config.StorageConnRetItv)
	if err != nil {
		return err
	}

	const query = `
        CREATE TABLE IF NOT EXISTS profile (
            user_id  TEXT
              CONSTRAINT profile_user_id_pk PRIMARY KEY,
            email    TEXT
              CONSTRAINT profile_email_unique UNIQUE,
            password TEXT,
            question TEXT,
            answer   TEXT
        );
	`

	_, err = r.conn.db.Exec(query)
	return err
}

func (r *ProfilesRepository) AddProfile(p *models.Profile) *errs.Error {
	const query = `
        INSERT INTO profile(user_id,email,password,question,answer) VALUES($1,$2,$3,$4,$5);
    `

	xRet := 3

start:
	p.UserId, _ = uuid.GenerateUUID()

	_, err := r.conn.db.Exec(query,
		p.UserId, p.Credentials.Email, p.Credentials.Password,
		p.SecretQuestion.Question, p.SecretQuestion.Answer,
	)
	if err != nil {
		if strings.Contains(err.Error(), "profile_user_id_pk") {
			if xRet == 0 {
				return errs.NewInternalError(err.Error())
			}
			xRet--
			goto start
		}
		if strings.Contains(err.Error(), "profile_email_unique") {
			return errs.NewConflictError("email duplicate")
		}
		return errs.NewInternalError(err.Error())
	}

	return nil
}

func (r *ProfilesRepository) GetProfileById(p *models.Profile) *errs.Error {
	panic("implement me")
}

func (r *ProfilesRepository) DeleteProfileById(p *models.Profile) *errs.Error {
	panic("implement me")
}
