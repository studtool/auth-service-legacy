package models

//go:generate easyjson

import (
	"github.com/studtool/common/types"
)

type Profile struct {
	UserId     string
	IsVerified bool

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
