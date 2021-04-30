package models

import (
	"encoding/json"
	"time"
)

// Date ...
type Date time.Time

// MarshalJSON ...
func (date Date) MarshalJSON() ([]byte, error) {
	formattedDate := time.Time(date).Format("2006-01-02")
	return json.Marshal(formattedDate)
}
