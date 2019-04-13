package types

import (
	"time"
)

const (
	DateLayout = "2006-01-02"
)

type Date time.Time

func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Date) UnmarshalJSON(b []byte) error {
	return t.Parse(string(b))
}

func (t Date) String() string {
	return time.Time(t).Format(DateLayout)
}

func (t *Date) Parse(s string) error {
	tm, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}

	*((*time.Time)(t)) = tm
	return nil
}
