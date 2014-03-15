package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"time"
)

type Bounce struct {
	CreatedAt string `json:"created_at"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	Error     string `json:"error"`
}

type bounceEnvelope struct {
	TotalCount int      `json:"total_count"`
	Items      []Bounce `json:"items"`
}

type singleBounceEnvelope struct {
	Bounce Bounce `json:"bounce"`
}

func (i Bounce) GetCreatedAt() (t time.Time, err error) {
	return parseMailgunTime(i.CreatedAt)
}

func (m *mailgunImpl) GetBounces(limit, skip int) (int, []Bounce, error) {
	r := simplehttp.NewHTTPRequest(generateApiUrl(m, bouncesEndpoint))
	if limit != -1 {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.AddParameter("skip", strconv.Itoa(skip))
	}

	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response bounceEnvelope
	err := r.GetResponseFromJSON(&response)
	if err != nil {
		return -1, nil, err
	}

	return response.TotalCount, response.Items, nil
}

func (m *mailgunImpl) GetSingleBounce(address string) (Bounce, error) {
	r := simplehttp.NewHTTPRequest(generateApiUrl(m, bouncesEndpoint) + "/" + address)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response singleBounceEnvelope
	err := r.GetResponseFromJSON(&response)
	if err != nil {
		return Bounce{}, err
	}

	return response.Bounce, nil
}

func (m *mailgunImpl) AddBounce(address, code, error string) error {
	r := simplehttp.NewHTTPRequest(generateApiUrl(m, bouncesEndpoint))
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	payload := simplehttp.NewUrlEncodedPayload()
	payload.AddValue("address", address)
	if code != "" {
		payload.AddValue("code", code)
	}
	if error != "" {
		payload.AddValue("error", error)
	}
	_, err := r.MakePostRequest(payload)
	return err
}

func (m *mailgunImpl) DeleteBounce(address string) error {
	r := simplehttp.NewHTTPRequest(generateApiUrl(m, bouncesEndpoint) + "/" + address)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeDeleteRequest()
	return err
}
