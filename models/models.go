package models

//go:generate easyjson

import (
	"github.com/studtool/auth-service/types"
)

//easyjson:json
type Profile struct {
	UserId string `json:"userId"`

	Credentials
}

//easyjson:json
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type Session struct {
	UserId       string         `json:"userId"`
	AuthToken    string         `json:"authToken"`
	RefreshToken string         `json:"refreshToken"`
	ExpireTime   types.DateTime `json:"expireTime"`
}
