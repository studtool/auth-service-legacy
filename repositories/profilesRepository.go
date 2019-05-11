package repositories

import (
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

type ProfilesRepository interface {
	AddProfile(p *models.Profile) *errs.Error
	SetProfileVerified(p *models.ProfileInfo) *errs.Error
	FindUserIdByCredentials(cr *models.Credentials) (string, *errs.Error)
	UpdateEmail(userId, email string) *errs.Error
	UpdatePassword(userId, password string) *errs.Error
	DeleteProfileById(userId string) *errs.Error
}
