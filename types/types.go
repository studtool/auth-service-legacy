package types

import (
	"time"
)

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format(time.RFC3339)), nil
}
