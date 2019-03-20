package models

import (
	"auth-service/errs"
	"regexp"
	"time"
)

//go:generate easyjson

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
	UserId       string    `json:"userId"`
	AuthToken    string    `json:"authToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpireTime   time.Time `json:"expireTime"`
}

type ProfileValidator struct {
	emailRegexp *regexp.Regexp
}

func NewProfileValidator() *ProfileValidator {
	return &ProfileValidator{
		emailRegexp: regexp.MustCompile(`^.+@.+$`),
	}
}

func (v *ProfileValidator) ValidateOnCreate(p *Profile) *errs.Error {
	if !v.emailRegexp.MatchString(p.Credentials.Email) {
		return errs.NewInvalidFormatError("invalid email")
	}
	if len(p.Credentials.Password) < 4 {
		return errs.NewInvalidFormatError("short password")
	}
	if len(p.SecretQuestion.Question) == 0 {
		return errs.NewInvalidFormatError("no secret question")
	}
	if len(p.SecretQuestion.Answer) == 0 {
		return errs.NewInvalidFormatError("no answer to secret question")
	}
	return nil
}
