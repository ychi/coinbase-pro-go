package common

import (
	"strings"
	"time"
)

type Time time.Time

func (t *Time) Time() time.Time{
	return time.Time(*t)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var parsed time.Time
	var err error
	if string(b) == "null" {
		*t = Time(time.Time{})
		return nil
	}

	layouts := []string{
		"2006-01-02 15:04:05+00",
		"2006-01-02T15:04:05.999999Z",
		"2006-01-02 15:04:05.999999",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05.999999+00",
	}

	for _, layout := range layouts {
		parsed, err = time.Parse(layout,
			strings.Replace(string(b), "\"", "", -1))
		if err != nil {
			continue
		}
		break
	}

	if parsed.IsZero() {
		return err
	}

	*t = Time(parsed)
	return nil
}

