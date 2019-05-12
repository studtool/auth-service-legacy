package utils

import (
	"time"

	"github.com/studtool/common/errs"
	"github.com/studtool/common/types"

	"github.com/studtool/auth-service/config"
)

type ExpireTimeCalculator struct {
	expiredErr  *errs.Error
	validPeriod time.Duration
}

func NewExpireTimeCalculator(err *errs.Error) *ExpireTimeCalculator {
	return &ExpireTimeCalculator{
		expiredErr:  err,
		validPeriod: config.JwtValidPeriod.Value(),
	}
}

func (c *ExpireTimeCalculator) Calculate() types.DateTime {
	return types.DateTime(time.Now().Add(c.validPeriod))
}

func (c *ExpireTimeCalculator) Check(t types.DateTime) *errs.Error {
	if time.Time(t).Before(time.Now()) {
		return c.expiredErr
	}
	return nil
}
