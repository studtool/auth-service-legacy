package repositories

import (
	"auth-service/models"
	"github.com/studtool/common/errs"
)

type ProfilesRepository interface {
	AddProfile(p *models.Profile) *errs.Error
	FindProfileByCredentials(p *models.Profile) *errs.Error
	UpdateCredentials(c *models.Credentials) *errs.Error
	DeleteProfileById(p *models.Profile) *errs.Error
}
