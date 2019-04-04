package repositories

import (
	"auth-service/models"
	"github.com/studtool/common/errs"
)

type SessionsRepository interface {
	AddSession(credentials *models.Credentials, session *models.Session) *errs.Error
	UpdateSessionByRefreshToken(session *models.Session) *errs.Error
	DeleteSessionByRefreshToken(token string) *errs.Error
}
