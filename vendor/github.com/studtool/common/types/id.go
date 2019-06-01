package types

import (
	"github.com/google/uuid"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
)

//go:generate msgp
//go:generate easyjson

type ID string

func MakeID() (ID, *errs.Error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return consts.EmptyString, errs.New(err)
	}
	return ID(id.String()), nil
}
