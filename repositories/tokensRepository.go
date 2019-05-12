package repositories

import (
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

type TokensRepository interface {
	SetToken(token *models.Token) *errs.Error
	GetToken(token *models.Token) *errs.Error
	DeleteToken(token *models.Token) *errs.Error
}
