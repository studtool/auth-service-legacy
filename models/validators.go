package models

import (
	"auth-service/errs"
	"fmt"
	"regexp"
)

type ProfileValidator struct {
	CredentialsValidator
	SecretQuestionValidator
}

func NewProfileValidator() *ProfileValidator {
	return &ProfileValidator{
		CredentialsValidator:    *NewCredentialsValidator(),
		SecretQuestionValidator: *NewSecretQuestionValidator(),
	}
}

func (v *ProfileValidator) ValidateOnCreate(p *Profile) *errs.Error {
	if err := v.CredentialsValidator.ValidateOnCreate(&p.Credentials); err != nil {
		return err
	}
	if err := v.SecretQuestionValidator.ValidateOnCreate(&p.SecretQuestion); err != nil {
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

const (
	questionPattern = `^(\w| |?)+$`
	answerPattern   = `^(\w| )+$`
)

type SecretQuestionValidator struct {
	questionRegexp *regexp.Regexp
	answerRegexp   *regexp.Regexp

	questionErr *errs.Error
	answerErr   *errs.Error
}

func NewSecretQuestionValidator() *SecretQuestionValidator {
	return &SecretQuestionValidator{
		questionRegexp: regexp.MustCompile(questionPattern),
		answerRegexp:   regexp.MustCompile(answerPattern),

		questionErr: errs.NewInvalidFormatError(
			fmt.Sprintf("invalid secret question. pattern: %s", questionPattern),
		),
		answerErr: errs.NewInvalidFormatError(
			fmt.Sprintf("invalid secret question answer: pattern: %s", answerPattern),
		),
	}
}

func (v *SecretQuestionValidator) ValidateOnCreate(obj *SecretQuestion) *errs.Error {
	if err := v.validateQuestion(obj.Question); err != nil {
		return err
	}
	if err := v.validateAnswer(obj.Answer); err != nil {
		return err
	}
	return nil
}

func (v *SecretQuestionValidator) ValidateOnUpdate(obj *SecretQuestion) *errs.Error {
	if obj.Question != "" {
		if err := v.validateQuestion(obj.Question); err != nil {
			return err
		}
	}

	if obj.Answer != "" {
		if err := v.validateAnswer(obj.Answer); err != nil {
			return err
		}
	}

	return nil
}

func (v *SecretQuestionValidator) validateQuestion(value string) *errs.Error {
	if !v.questionRegexp.MatchString(value) {
		return v.questionErr
	}
	return nil
}

func (v *SecretQuestionValidator) validateAnswer(value string) *errs.Error {
	if !v.answerRegexp.MatchString(value) {
		return v.questionErr
	}
	return nil
}
