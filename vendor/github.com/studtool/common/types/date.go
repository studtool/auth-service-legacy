package types

import (
	"encoding/json"
	"time"
)

const (
	DateLayout = "2006-01-02"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	return d.Parse(str)
}

func (d Date) String() string {
	return time.Time(d).Format(DateLayout)
}

func (d *Date) Parse(s string) error {
	tm, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}

	*((*time.Time)(d)) = tm
	return nil
}
