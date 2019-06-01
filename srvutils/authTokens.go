package srvutils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
	"github.com/studtool/common/types"

	"github.com/studtool/auth-service/config"
)

type AuthTokenAttributes struct {
	UserID  string
	ExpTime types.DateTime
}

type jwtClaims struct {
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

func (m *AuthTokenManager) CreateToken(attr *AuthTokenAttributes) (string, *errs.Error) {
	jwtClaims := jwtClaims{
		UserId:  attr.UserID,
		ExpTime: attr.ExpTime.String(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  jwtClaims.UserId,
		"expTime": jwtClaims.ExpTime,
	})

	if s, err := t.SignedString(m.key); err != nil {
		return consts.EmptyString, errs.New(err)
	} else {
		return s, nil
	}
}

func (m *AuthTokenManager) ParseToken(token string) (*AuthTokenAttributes, *errs.Error) {
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

	jwtClaims := &jwtClaims{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: jwtClaims, TagName: "mapstructure",
	})
	if err != nil {
		return nil, errs.New(err)
	}

	if err := decoder.Decode(claims); err != nil {
		return nil, m.err
	}

	attr := &AuthTokenAttributes{
		UserID: jwtClaims.UserId,
	}
	if err := attr.ExpTime.Parse(jwtClaims.ExpTime); err != nil {
		return nil, errs.New(err)
	}

	return attr, nil
}
