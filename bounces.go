package mailgun

import (
	"time"
)

type BounceItem struct {
	CreatedAt time.Time `json:"created_at"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	Error     string    `json:"error"`
}

type Bounces struct {
	TotalCount int          `json:"total_count"`
	Items      []BounceItem `json:"items"`
}
