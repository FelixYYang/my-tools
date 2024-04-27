package types

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.DateOnly)+len(`""`))
	b = append(b, '"')
	t := time.Time(d)
	b = t.AppendFormat(b, time.DateOnly)
	b = append(b, '"')
	return b, nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}
	t, err := time.Parse(time.DateOnly, string(b[1:len(b)-1]))
	*d = Date(t)
	return err
}

func (d *Date) Scan(src any) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(src)
	*d = Date(nullTime.Time)
	return
}

func (d Date) Value() (driver.Value, error) {
	y, m, day := time.Time(d).Date()
	return time.Date(y, m, day, 0, 0, 0, 0, time.Time(d).Location()), nil
}

type DateTime time.Time

func (d DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.DateTime)+len(`""`))
	b = append(b, '"')
	t := time.Time(d)
	b = t.AppendFormat(b, time.DateTime)
	b = append(b, '"')
	return b, nil
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	return (*time.Time)(d).UnmarshalJSON(b)
}
