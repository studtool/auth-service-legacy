package models

//go:generate easyjson

import (
	"auth-service/types"
)

//easyjson:json
type Profile struct {
	UserId         string         `json:"userId"`
	Credentials    Credentials    `json:"credentials"`
	SecretQuestion SecretQuestion `json:"secretQuestion"`
}

//easyjson:json
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type SecretQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

//easyjson:json
type Session struct {
	UserId       string     `json:"userId"`
	AuthToken    string     `json:"authToken"`
	RefreshToken string     `json:"refreshToken"`
	ExpireTime   types.Time `json:"expireTime"`
}
