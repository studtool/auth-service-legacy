package models

import (
	"auth-service/errs"
	"fmt"
	"regexp"
)

type ProfileValidator struct {
	CredentialsValidator
}

func NewProfileValidator() *ProfileValidator {
	return &ProfileValidator{
		CredentialsValidator: *NewCredentialsValidator(),
	}
}

func (v *ProfileValidator) ValidateOnCreate(p *Profile) *errs.Error {
	if err := v.CredentialsValidator.ValidateOnCreate(&p.Credentials); err != nil {
		return err
	}
	return nil
}

const (
	emailPattern      = `^.+@.+$`
	minPasswordLength = 4
)

type CredentialsValidator struct {
	emailRegexp *regexp.Regexp

	emailErr    *errs.Error
	passwordErr *errs.Error
}

func NewCredentialsValidator() *CredentialsValidator {
	return &CredentialsValidator{
		emailRegexp: regexp.MustCompile(emailPattern),

		emailErr: errs.NewInvalidFormatError(
			fmt.Sprintf("invalid email. pattern: " + emailPattern),
		),
		passwordErr: errs.NewInvalidFormatError(
			fmt.Sprintf("invalid password. min length: %d", minPasswordLength),
		),
	}
}

func (v *CredentialsValidator) ValidateOnCreate(obj *Credentials) *errs.Error {
	if err := v.validateEmail(obj.Email); err != nil {
		return err
	}
	if err := v.validatePassword(obj.Password); err != nil {
		return err
	}
	return nil
}

func (v *CredentialsValidator) ValidateOnUpdate(obj *Credentials) *errs.Error {
	if obj.Email != "" {
		if err := v.validateEmail(obj.Email); err != nil {
			return err
		}
	}

	if obj.Password != "" {
		if err := v.validatePassword(obj.Password); err != nil {
			return err
		}
	}

	return nil
}

func (v *CredentialsValidator) validateEmail(value string) *errs.Error {
	if !v.emailRegexp.MatchString(value) {
		return v.emailErr
	}
	return nil
}

func (v *CredentialsValidator) validatePassword(value string) *errs.Error {
	if len(value) < minPasswordLength {
		return v.passwordErr
	}
	return nil
}
