package repositories

import (
	"auth-service/errs"
	"auth-service/models"
)

type SessionsRepository interface {
	AddSession(session *models.Session) *errs.Error
	UpdateSessionByRefreshToken(session *models.Session) *errs.Error
	DeleteSessionByRefreshToken(token string) *errs.Error
}
