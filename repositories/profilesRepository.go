package repositories

import (
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

type ProfilesRepository interface {
	AddProfile(p *models.Profile) *errs.Error
	SetProfileVerified(p *models.ProfileInfo) *errs.Error
	FindUserIdByCredentials(p *models.Profile) *errs.Error
	UpdateCredentials(c *models.Credentials) *errs.Error
	DeleteProfileById(p *models.Profile) *errs.Error
}
