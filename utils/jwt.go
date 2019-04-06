package utils

//go:generate easyjson

import (
	"auth-service/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/mailru/easyjson"
	"github.com/studtool/common/errs"
)

//easyjson:json
type JwtClaims struct {
	UserId  string `json:"userId"`
	ExpTime string `json:"expTime"`
}

type JwtManager struct {
	key []byte
	err *errs.Error
}

func NewJwtManager() *JwtManager {
	return &JwtManager{
		key: []byte(config.JwtKey.Value()),
		err: errs.NewNotAuthorizedError("invalid token"),
	}
}

func (m *JwtManager) CreateToken(c *JwtClaims) (string, *errs.Error) {
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

func (m *JwtManager) ParseToken(token string) (*JwtClaims, *errs.Error) {
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
