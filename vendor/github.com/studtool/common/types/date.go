package types

import (
	"time"
)

const (
	DateLayout = "2006-01-02"
)

type Date time.Time

func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(DateLayout)), nil
}

func (t *Date) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(DateLayout, string(b))
	if err != nil {
		return err
	}

	*((*time.Time)(t)) = tm
	return nil
}
