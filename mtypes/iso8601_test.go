package mtypes

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
		wantErr bool
	}{
		{
			name:    "ISO8601",
			s:       `{"created_at":"2011/10/13T18:02:00"}`,
			wantErr: false,
			want:    Req{CreatedAt: ISO8601Time(time.Date(2011, 10, 13, 18, 2, 0, 0, time.UTC))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req Req
			err := json.Unmarshal([]byte(tt.s), &req)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}

			require.True(t, time.Time(tt.want.CreatedAt).Equal(time.Time(req.CreatedAt)),
				fmt.Sprintf("want: %s; got: %s", tt.want.CreatedAt, req.CreatedAt))
		})
	}
}
