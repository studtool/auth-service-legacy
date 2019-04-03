package repositories

import (
	"auth-service/errs"
	"auth-service/models"
)

type ProfilesRepository interface {
	AddProfile(p *models.Profile) *errs.Error
	GetProfileByCredentials(p *models.Profile) *errs.Error
	UpdateCredentials(c *models.Credentials) *errs.Error
	DeleteProfileById(p *models.Profile) *errs.Error
}
