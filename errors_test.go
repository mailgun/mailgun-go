package mailgun

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newError(t *testing.T) {
	const (
		method = "GET"
		uri    = "/v1/foo"
	)
	expected := []int{200, 202, 204}

	tests := []struct {
		name         string
		method       string
		uri          string
		expected     []int
		httpResponse httpResponse
		wantErr      error
	}{
		{
			name:         "400",
			method:       method,
			uri:          uri,
			expected:     expected,
			httpResponse: httpResponse{Code: 400, Data: []byte("Bad Request")},
			wantErr: &UnexpectedResponseError{
				Expected: expected,
				Actual:   400,
				Method:   method,
				URL:      uri,
				Data:     []byte("Bad Request"),
			},
		},

		{
			name:         "429",
			method:       method,
			uri:          uri,
			expected:     expected,
			httpResponse: httpResponse{Code: 429, Data: []byte("Too Many Requests")},
			wantErr: &RateLimitedError{
				Err: &UnexpectedResponseError{
					Expected: expected,
					Actual:   429,
					Method:   method,
					URL:      uri,
					Data:     []byte("Too Many Requests"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := newError(tt.method, tt.uri, tt.expected, &tt.httpResponse)
			if tt.wantErr != nil {
				assert.EqualError(t, gotErr, tt.wantErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestRateLimitedError(t *testing.T) {
	err := newError("GET", "/v1/foo", []int{200, 201}, &httpResponse{Code: 429, Data: []byte("Too Many Requests")})

	t.Run(".Error()", func(t *testing.T) {
		wantErr := errors.New(`RateLimitedError: UnexpectedResponseError Method=GET URL=/v1/foo ExpectedOneOf=[]int{200, 201} Got=429 Error: Too Many Requests`)

		assert.EqualError(t, err, wantErr.Error())
	})

	t.Run("GetStatusFromErr()", func(t *testing.T) {
		status := GetStatusFromErr(err)

		assert.Equal(t, http.StatusTooManyRequests, status)
	})
}
