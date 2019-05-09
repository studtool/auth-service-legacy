package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/config"
)

type JwtClaims struct {
	UserId  string `mapstructure:"userId"`
	ExpTime string `mapstructure:"expTime"`
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
		return consts.EmptyString, errs.New(err)
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
		return nil, errs.New(err)
	}

	if err := decoder.Decode(claims); err != nil {
		return nil, m.err
	}

	return jwtClaims, nil
}
