package utils

import (
	"fmt"

	"github.com/studtool/common/errs"
	"github.com/studtool/common/utils"
)

const (
	refreshTokenRandLen = 128
)

type RefreshTokenAttributes struct {
	UserID string
}

type RefreshTokenManager struct{}

func NewRefreshTokenManager() *RefreshTokenManager {
	return &RefreshTokenManager{}
}

func (m *RefreshTokenManager) CreateToken(attr *RefreshTokenAttributes) (string, *errs.Error) {
	t := fmt.Sprintf("%s%s",
		attr.UserID, utils.RandString(refreshTokenRandLen),
	)
	return t, nil
}
