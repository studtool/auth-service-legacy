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
	return t.Parse(string(b))
}

func (t DateTime) String() string {
	return time.Time(t).Format(DateTimeLayout)
}

func (t *DateTime) Parse(s string) error {
	tm, err := time.Parse(DateTimeLayout, s)
	if err != nil {
		return err
	}

	*((*time.Time)(t)) = tm
	return nil
}
