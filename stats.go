package mailgun

import (
	"context"
	"strconv"
	"time"
)

// Stats on accepted messages
type Accepted struct {
	Incoming int `json:"incoming"`
	Outgoing int `json:"outgoing"`
	Total    int `json:"total"`
}

// Stats on delivered messages
type Delivered struct {
	Smtp  int `json:"smtp"`
	Http  int `json:"http"`
	Total int `json:"total"`
}

// Stats on temporary failures
type Temporary struct {
	Espblock int `json:"espblock"`
}

// Stats on permanent failures
type Permanent struct {
	SuppressBounce      int `json:"suppress-bounce"`
	SuppressUnsubscribe int `json:"suppress-unsubscribe"`
	SuppressComplaint   int `json:"suppress-complaint"`
	Bounce              int `json:"bounce"`
	DelayedBounce       int `json:"delayed-bounce"`
	Total               int `json:"total"`
}

// Stats on failed messages
type Failed struct {
	Temporary Temporary `json:"temporary"`
	Permanent Permanent `json:"permanent"`
}

// Total stats for messages
type Total struct {
	Total int `json:"total"`
}

// Stats as returned by `GetStats()`
type Stats struct {
	Time         string    `json:"time"`
	Accepted     Accepted  `json:"accepted"`
	Delivered    Delivered `json:"delivered"`
	Failed       Failed    `json:"failed"`
	Stored       Total     `json:"stored"`
	Opened       Total     `json:"opened"`
	Clicked      Total     `json:"clicked"`
	Unsubscribed Total     `json:"unsubscribed"`
	Complained   Total     `json:"complained"`
}

type statsTotalResponse struct {
	End        string  `json:"end"`
	Resolution string  `json:"resolution"`
	Start      string  `json:"start"`
	Stats      []Stats `json:"stats"`
}

// Used by GetStats() to specify the resolution stats are for
type Resolution string

// Indicate which resolution a stat response for request is for
const (
	ResolutionHour  = Resolution("hour")
	ResolutionDay   = Resolution("day")
	ResolutionMonth = Resolution("month")
)

// Options for GetStats()
type GetStatOptions struct {
	Resolution Resolution
	Duration   string
	Start      time.Time
	End        time.Time
}

// GetStats returns total stats for a given domain for the specified time period
func (mg *MailgunImpl) GetStats(ctx context.Context, events []string, opts *GetStatOptions) ([]Stats, error) {
	r := newHTTPRequest(generateApiUrl(mg, statsTotalEndpoint))

	if opts != nil {
		if !opts.Start.IsZero() {
			r.addParameter("start", strconv.Itoa(int(opts.Start.Unix())))
		}
		if !opts.End.IsZero() {
			r.addParameter("end", strconv.Itoa(int(opts.End.Unix())))
		}
		if opts.Resolution != "" {
			r.addParameter("resolution", string(opts.Resolution))
		}
		if opts.Duration != "" {
			r.addParameter("duration", opts.Duration)
		}
	}

	for _, e := range events {
		r.addParameter("event", e)
	}

	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var res statsTotalResponse
	err := getResponseFromJSON(ctx, r, &res)
	if err != nil {
		return nil, err
	} else {
		return res.Stats, nil
	}
}
