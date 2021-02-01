package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type DecimalStr string

func (d *DecimalStr) UnmarshalJSON(b []byte) error {
	var v interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case float64:
		*d = DecimalStr(fmt.Sprintf("%f", v))
	case int:
		*d = DecimalStr(strconv.Itoa(v))
	case string:
		*d = DecimalStr(v)
	default:
		return errors.New("Not float, int or string")
	}

	return nil
}
