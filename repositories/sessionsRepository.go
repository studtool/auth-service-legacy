package repositories

import (
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

type SessionsRepository interface {
	AddSession(session *models.Session) *errs.Error
	FindSession(session *models.Session) *errs.Error
	DeleteSessionByRefreshToken(token string) *errs.Error
	DeleteAllSessionsByRefreshToken(token string) *errs.Error
}
