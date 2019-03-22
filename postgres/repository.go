package postgres

import (
	"auth-service/beans"
	"auth-service/config"
	"auth-service/errs"
	"auth-service/models"
	"database/sql"
	"fmt"
	"github.com/hashicorp/go-uuid"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

const (
	InitRetryPeriod = 2 * time.Second
)

type Repository struct {
	connStr string
	db      *sql.DB
}

func NewRepository() *Repository {
	return &Repository{
		connStr: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.StorageUser, config.StoragePassword,
			config.StorageHost, config.StoragePort,
			config.StorageDB, config.StorageSSL,
		),
	}
}

func (r *Repository) Open() {
	var err error

	r.db, err = sql.Open("postgres", r.connStr)
	if err != nil {
		panic(err)
	}

	beans.Logger.Infof(
		"postgres repository on %s:%s",
		config.StorageHost, config.StoragePort,
	)
}

func (r *Repository) Close() {
	if err := r.db.Close(); err != nil {
		panic(err)
	}
	beans.Logger.Infof("postgres repository connection closed")
}

func (r *Repository) Init() {
	if !config.ShouldInitStorage {
		return
	}

	xRet := 1
	for r.db.Ping() != nil {
		beans.Logger.Infof("waiting for database initialization: retry #%d", xRet)
		xRet++

		time.Sleep(InitRetryPeriod)
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

	_, err := r.db.Exec(query)
	if err != nil {
		panic(err)
	}

	beans.Logger.Debugf("db initialized: \n%s", query)
}

func (r *Repository) CreateProfile(p *models.Profile) *errs.Error {
	const query = `
        INSERT INTO profile(user_id,email,password,question,answer) VALUES($1,$2,$3,$4,$5);
    `

	xRet := 3

start:
	p.UserId, _ = uuid.GenerateUUID()

	_, err := r.db.Exec(query,
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

func (r *Repository) CreateSession(p *models.Session) *errs.Error {
	return nil //TODO
}
