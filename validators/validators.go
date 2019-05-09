package validators

import (
	"fmt"
	"regexp"

	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/models"
)

const (
	emailPattern      = `^.+@.+$`
	minPasswordLength = 4
)

type CredentialsValidator struct {
	emailValidator    *EmailValidator
	passwordValidator *PasswordValidator
}

func NewCredentialsValidator() *CredentialsValidator {
	return &CredentialsValidator{
		emailValidator:    NewEmailValidator(),
		passwordValidator: NewPasswordValidator(),
	}
}

func (v *CredentialsValidator) Validate(obj *models.Credentials) *errs.Error {
	if err := v.emailValidator.Validate(obj.Email); err != nil {
		return err
	}
	if err := v.passwordValidator.Validate(obj.Password); err != nil {
		return err
	}
	return nil
}

type EmailValidator struct {
	re  *regexp.Regexp
	err *errs.Error
}

func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		re: regexp.MustCompile(emailPattern),
		err: errs.NewInvalidFormatError(
			fmt.Sprintf("invalid email. pattern: %s", emailPattern),
		),
	}
}

func (v *EmailValidator) Validate(email string) *errs.Error {
	if !v.re.MatchString(email) {
		return v.err
	}
	return nil
}

type PasswordValidator struct {
	err *errs.Error
}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		err: errs.NewInvalidFormatError(
			fmt.Sprintf("invalid password. min length: %d", minPasswordLength),
		),
	}
}

func (v *PasswordValidator) Validate(password string) *errs.Error {
	if len(password) < minPasswordLength {
		return v.err
	}
	return nil
}
