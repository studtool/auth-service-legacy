package utils

//go:generate easyjson

import (
	"auth-service/config"
	"auth-service/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/mailru/easyjson"
	"github.com/studtool/common/errs"
)

//easyjson:json
type JwtClaims struct {
	UserId  string         `json:"userId"`
	ExpTime types.DateTime `json:"expTime"`
}

type AuthTokenManager struct {
	key []byte
	err *errs.Error
}

func NewAuthTokenManager() *AuthTokenManager {
	return &AuthTokenManager{
		key: []byte(config.JwtKey.Value()),
		err: errs.NewNotAuthorizedError("invalid token"),
	}
}

func (m *AuthTokenManager) CreateToken(c *JwtClaims) (string, *errs.Error) {
	d, _ := easyjson.Marshal(c)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": d,
	})

	if s, err := t.SignedString(m.key); err != nil {
		return "", errs.NewInternalError(err.Error())
	} else {
		return s, nil
	}
}

func (m *AuthTokenManager) ParseToken(token string) (*JwtClaims, *errs.Error) {
	t, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, m.err
		}
		return m.key, nil
	})

	if err != nil || !t.Valid {
		return nil, m.err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, m.err
	}

	data, ok := claims["data"].([]byte)
	if !ok {
		return nil, m.err
	}

	jwtClaims := &JwtClaims{}
	if err := easyjson.Unmarshal(data, jwtClaims); err != nil {
		return nil, m.err
	}

	return jwtClaims, nil
}
