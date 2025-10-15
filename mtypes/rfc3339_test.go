package mtypes

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalRFC3339Time(t *testing.T) {
	type Req struct {
		CreatedAt RFC3339Time `json:"created_at"`
	}

	tests := []struct {
		name    string
		s       string
		want    Req
		wantErr bool
	}{
		{
			name:    "RFC1123",
			s:       `{"created_at":"2025-01-01T01:23:45Z"}`,
			wantErr: false,
			want:    Req{CreatedAt: RFC3339Time(time.Date(2025, 1, 1, 1, 23, 45, 0, time.UTC))},
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
