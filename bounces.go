package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"time"
)

// Bounce aggregates data relating to undeliverable messages to a specific intended recipient,
// identified by Address.
// Code provides the SMTP error code causing the bounce,
// while Error provides a human readable reason why.
// CreatedAt provides the time at which Mailgun detected the bounce.
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

// GetCreatedAt parses the textual, RFC-822 timestamp into a standard Go-compatible
// Time structure.
func (i Bounce) GetCreatedAt() (t time.Time, err error) {
	return parseMailgunTime(i.CreatedAt)
}

// GetBounces returns a complete set of bounces logged against the sender's domain, if any.
func (m *mailgunImpl) GetBounces(limit, skip int) (int, []Bounce, error) {
	r := simplehttp.NewGetRequest(generateApiUrl(m, bouncesEndpoint))
	if limit != -1 {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.AddParameter("skip", strconv.Itoa(skip))
	}

	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response bounceEnvelope
	err := r.MakeJSONRequest(&response)
	if err != nil {
		return -1, nil, err
	}

	return response.TotalCount, response.Items, nil
}

// GetSingleBounce retrieves a single bounce record, if any exist, for the given recipient address.
// If none exist, the returned Bounce instance will be empty.
func (m *mailgunImpl) GetSingleBounce(address string) (Bounce, error) {
	r := simplehttp.NewGetRequest(generateApiUrl(m, bouncesEndpoint) + "/" + address)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response singleBounceEnvelope
	err := r.MakeJSONRequest(&response)
	if err != nil {
		return Bounce{}, err
	}

	return response.Bounce, nil
}

// AddBounce files a bounce report.
// Address identifies the intended recipient of the message that bounced.
// Code corresponds to the numeric response given by the e-mail server which rejected the message.
// Error providees the corresponding human readable reason for the problem.
// For example,
// here's how the these two fields relate.
// Suppose the SMTP server responds with an error, as below.
// Then, . . .
//
//      550  Requested action not taken: mailbox unavailable
//     \___/\_______________________________________________/
//       |                         |
//       +-- Code                  +-- Error
//
// Note that both code and error exist as strings, even though
// code will report as a number.
func (m *mailgunImpl) AddBounce(address, code, error string) error {
	r := simplehttp.NewPostRequest(generateApiUrl(m, bouncesEndpoint))

	r.AddFormValue("address", address)
	if code != "" {
		r.AddFormValue("code", code)
	}
	if error != "" {
		r.AddFormValue("error", error)
	}
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}

// DeleteBounce removes all bounces associted with the provided e-mail address.
func (m *mailgunImpl) DeleteBounce(address string) error {
	r := simplehttp.NewDeleteRequest(generateApiUrl(m, bouncesEndpoint) + "/" + address)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}
