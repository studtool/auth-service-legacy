package models

//go:generate easyjson

import (
	"github.com/studtool/common/types"
)

type Profile struct {
	Credentials
	ProfileInfo
}

//easyjson:json
type ProfileInfo struct {
	UserID     string `json:"userId"`
	IsVerified bool   `json:"isVerified"`
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

//easyjson:json
type Token struct {
	UserID string `json:"-"`
	Token  string `json:"token"`
}
