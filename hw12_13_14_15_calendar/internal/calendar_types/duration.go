package calendar_types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CalendarDuration time.Duration

func (calendarDuration *CalendarDuration) UnmarshalJSON(bytes []byte) error {
	strs := strings.Trim(string(bytes), `"`)
	duration, err := time.ParseDuration(strs)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	*calendarDuration = CalendarDuration(duration)
	return nil
}

func (calendarDuration CalendarDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Duration(calendarDuration).String())), nil
}

func (calendarDuration CalendarDuration) Value() (driver.Value, error) {
	return time.Duration(calendarDuration).String(), nil
}

func (calendarDuration *CalendarDuration) Scan(src interface{}) error {
	var str string
	switch value := src.(type) {
	case string:
		str = value
	case []byte:
		str = string(value)
	default:
		return fmt.Errorf("can't scan value of type %T", src)
	}

	timeDelta, err := time.Parse("15:04:05", str)
	if err != nil {
		//return fmt.Errorf("can't scan duration: %w", err)
		// Попытка извлечь "огромное число часов"
		parts := strings.Split(str, ":")
		if len(parts) == 3 {
			hours, err1 := strconv.ParseInt(parts[0], 10, 64)
			mins, err2 := strconv.ParseInt(parts[1], 10, 64)
			secs, err3 := strconv.ParseInt(parts[2], 10, 64)

			if err1 == nil && err2 == nil && err3 == nil {
				dur := time.Duration(hours)*time.Hour +
					time.Duration(mins)*time.Minute +
					time.Duration(secs)*time.Second
				*calendarDuration = CalendarDuration(dur)
				return nil
			}
			return fmt.Errorf("can't scan value of type %T", src)
		}
	}

	duration := time.Duration(timeDelta.Hour())*time.Hour +
		(time.Duration(timeDelta.Minute()) * time.Minute) +
		(time.Duration(timeDelta.Second()) * time.Second)

	*calendarDuration = CalendarDuration(duration)
	return nil
}
