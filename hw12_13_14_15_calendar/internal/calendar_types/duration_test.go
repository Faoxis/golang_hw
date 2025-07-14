package calendar_types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalendarDuration_JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    CalendarDuration
		expected string
	}{
		//nolint:typecheck
		{"zero duration", CalendarDuration(0), `"0s"`},
		//nolint:typecheck
		{"minutes and seconds", CalendarDuration(90 * time.Second), `"1m30s"`},
		//nolint:typecheck
		{"hours and minutes", CalendarDuration(2*time.Hour + 15*time.Minute), `"2h15m0s"`},
		{"negative duration", CalendarDuration(-10 * time.Second), `"-10s"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestCalendarDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonIn   string
		expected CalendarDuration
		hasError bool
	}{
		{"valid duration", `"1h30m"`, CalendarDuration(90 * time.Minute), false},
		{"empty string", `""`, CalendarDuration(time.Duration(0)), true},
		{"invalid format", `"abc"`, CalendarDuration(time.Duration(0)), true},
		{"negative duration", `"-5s"`, CalendarDuration(-5 * time.Second), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d CalendarDuration
			err := json.Unmarshal([]byte(tt.jsonIn), &d)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, d)
			}
		})
	}
}

func TestCalendarDuration_ValueAndScan(t *testing.T) {
	t.Run("Value returns string", func(t *testing.T) {
		d := CalendarDuration(45 * time.Second)
		val, err := d.Value()
		assert.NoError(t, err)
		assert.Equal(t, "45s", val)
	})

	scanTests := []struct {
		name     string
		input    interface{}
		expected CalendarDuration
		hasError bool
	}{
		{"from string", "1:0:0", CalendarDuration(time.Hour), false},
		{"from []byte", []byte("2:30:0"), CalendarDuration(2*time.Hour + 30*time.Minute), false},
		{"invalid type", 123, CalendarDuration(time.Duration(0)), true},
	}

	for _, tt := range scanTests {
		t.Run("Scan: "+tt.name, func(t *testing.T) {
			var d CalendarDuration
			err := d.Scan(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, d)
			}
		})
	}
}
