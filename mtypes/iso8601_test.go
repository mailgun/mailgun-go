package mtypes

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalISO8601Time(t *testing.T) {
	type Req struct {
		CreatedAt ISO8601Time `json:"created_at"`
	}

	tests := []struct {
		name    string
		s       string
		want    Req
		wantErr error
	}{
		{
			name:    "ISO8601",
			s:       `{"created_at":"2011-10-13T18:02:00"}`,
			want:    Req{CreatedAt: ISO8601Time{time.Date(2011, 10, 13, 18, 2, 0, 0, time.UTC)}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req Req
			err := json.Unmarshal([]byte(tt.s), &req)
			if tt.wantErr != nil {
				require.Contains(t, err.Error(), tt.wantErr.Error())
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.want.CreatedAt.Equal(req.CreatedAt.Time),
				fmt.Sprintf("want: %s; got: %s", tt.want.CreatedAt, req.CreatedAt))
		})
	}
}
