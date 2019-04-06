package utils

//go:generate easyjson

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"

	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/types"
)

//easyjson:json
type JwtClaims struct {
	UserId  string         `mapstructure:"userId"`
	ExpTime types.DateTime `mapstructure:"expTime"`
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
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  c.UserId,
		"expTime": c.ExpTime,
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

	jwtClaims := &JwtClaims{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: jwtClaims, TagName: "mapstructure",
	})
	if err != nil {
		panic(err)
	}

	if err := decoder.Decode(claims); err != nil {
		return nil, m.err
	}

	return jwtClaims, nil
}
