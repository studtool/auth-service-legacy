package types

import (
	"time"
)

const (
	DateTimeLayout = time.RFC3339
)

type DateTime time.Time

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(DateTimeLayout)), nil
}

func (t *DateTime) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(DateTimeLayout, string(b))
	if err != nil {
		return err
	}

	*((*time.Time)(t)) = tm
	return nil
}
