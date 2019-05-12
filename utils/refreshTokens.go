package utils

import (
	"github.com/google/uuid"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/beans"
)

const (
	uuidChainLen = 5
)

type RefreshTokenManager struct{}

func NewRefreshTokenManager() *RefreshTokenManager {
	return &RefreshTokenManager{}
}

func (m *RefreshTokenManager) CreateToken() (string, *errs.Error) {
	t := consts.EmptyString
	for i := 1; i <= uuidChainLen; i++ {
		t += m.generateUUID()
	}
	return t, nil
}

func (m *RefreshTokenManager) generateUUID() string {
	v, err := uuid.NewRandom()
	if err != nil { //TODO handle error
		beans.Logger().Error(err)
	}
	return v.String()
}
