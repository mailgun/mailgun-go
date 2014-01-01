package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"time"
)

type Stat struct {
	Event      string         `json:"event"`
	TotalCount int            `json:"total_count"`
	CreatedAt  string         `json:"created_at"`
	Id         string         `json:"id"`
	Tags       map[string]int `json:"tags"`
}

type statsEnvelope struct {
	TotalCount int    `json:"total_count"`
	Items      []Stat `json:"items"`
}

func (m *mailgunImpl) GetStats(limit int, skip int, startDate time.Time, event ...string) (int, []Stat, error) {
	r := simplehttp.NewGetRequest(generateApiUrl(m, statsEndpoint))

	if limit != -1 {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.AddParameter("skip", strconv.Itoa(skip))
	}

	r.AddParameter("start-date", startDate.Format("2006-02-01"))

	for _, e := range event {
		r.AddParameter("event", e)
	}
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var res statsEnvelope
	err := r.MakeJSONRequest(&res)
	if err != nil {
		return -1, nil, err
	} else {
		return res.TotalCount, res.Items, nil
	}
}

func (m *mailgunImpl) DeleteTag(tag string) error {
	r := simplehttp.NewDeleteRequest(generateApiUrl(m, deleteTagEndpoint) + "/" + tag)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}
