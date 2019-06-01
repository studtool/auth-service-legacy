package types

//go:generate easyjson

import (
	"encoding/json"
	"time"
)

const (
	DateTimeLayout = time.RFC3339
)

type DateTime time.Time

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	return dt.Parse(str)
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(DateTimeLayout)
}

func (dt *DateTime) Parse(s string) error {
	tm, err := time.Parse(DateTimeLayout, s)
	if err != nil {
		return err
	}

	*((*time.Time)(dt)) = tm
	return nil
}
