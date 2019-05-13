package repositories

import (
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

type ProfilesRepository interface {
	AddProfile(p *models.Profile) *errs.Error
	SetProfileVerified(p *models.ProfileInfo) *errs.Error
	FindProfile(p *models.Profile) *errs.Error
	FindVerifiedProfile(p *models.Profile) *errs.Error
	UpdateEmail(u *models.EmailUpdate) *errs.Error
	UpdatePassword(u *models.PasswordUpdate) *errs.Error
	DeleteProfileById(userId string) *errs.Error
}
