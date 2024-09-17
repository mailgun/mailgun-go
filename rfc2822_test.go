package mailgun

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
)

func TestUnmarshalRFC2822Time(t *testing.T) {
	type Req struct {
		CreatedAt RFC2822Time `json:"created_at"`
	}

	tests := []struct {
		name    string
		s       string
		want    Req
		wantErr bool
	}{
		{
			name:    "RFC1123",
			s:       `{"created_at":"Thu, 13 Oct 2011 18:02:00 GMT"}`,
			wantErr: false,
			want:    Req{CreatedAt: RFC2822Time(time.Date(2011, 10, 13, 18, 2, 0, 0, time.UTC))},
		},
		{
			name:    "RFC1123Z",
			s:       `{"created_at":"Thu, 13 Oct 2011 18:02:00 +0000"}`,
			wantErr: false,
			want:    Req{CreatedAt: RFC2822Time(time.Date(2011, 10, 13, 18, 2, 0, 0, time.UTC))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req Req
			err := json.Unmarshal([]byte(tt.s), &req)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}

			ensure.True(t, time.Time(tt.want.CreatedAt).Equal(time.Time(req.CreatedAt)),
				fmt.Sprintf("want: %s; got: %s", tt.want.CreatedAt, req.CreatedAt))
		})
	}
}
