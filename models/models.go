package models

//go:generate easyjson

import (
	"auth-service/types"
)

//easyjson:json
type Profile struct {
	UserId      string      `json:"userId"`
	Credentials Credentials `json:"credentials"`
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
