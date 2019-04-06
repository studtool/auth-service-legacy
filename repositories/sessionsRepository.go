package repositories

import (
	"auth-service/models"
	"github.com/studtool/common/errs"
)

type SessionsRepository interface {
	AddSession(session *models.Session) *errs.Error
	FindUserIdByRefreshToken(session *models.Session) *errs.Error
	UpdateSessionByRefreshToken(session *models.Session) *errs.Error
	DeleteSessionByRefreshToken(token string) *errs.Error
}
