package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		date    Date
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "with a non-zero time",
			date:    Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
			want:    []byte(`"2006-01-02"`),
			wantErr: assert.NoError,
		},
		{
			name:    "with a zero time",
			date:    Date(time.Time{}),
			want:    []byte(`"0001-01-01"`),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.date)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}