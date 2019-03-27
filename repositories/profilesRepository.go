package repositories

import (
	"auth-service/errs"
	"auth-service/models"
)

type ProfilesRepository interface {
	AddProfile(p *models.Profile) *errs.Error
	GetProfileById(p *models.Profile) *errs.Error
	DeleteProfileById(p *models.Profile) *errs.Error
}
