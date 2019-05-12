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
type EmailUpdate struct {
	UserID string `json:"-"`
	Email  string `json:"email"`
}

//easyjson:json
type PasswordUpdate struct {
	UserID   string `json:"-"`
	Password string `json:"password"`
}

//easyjson:json
type Session struct {
	SessionID    string         `json:"sessionId"`
	UserID       string         `json:"userId"`
	AuthToken    string         `json:"authToken"`
	RefreshToken string         `json:"refreshToken"`
	ExpireTime   types.DateTime `json:"expireTime"`
}

//easyjson:json
type Token struct {
	UserID string `json:"-"`
	Token  string `json:"token"`
}
