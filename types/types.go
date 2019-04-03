package types

import (
	"time"
)

type DateTime time.Time

func (t *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format(time.RFC3339)), nil
}
