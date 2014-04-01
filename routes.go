package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
)

// A Route structure contains information on a configured or to-be-configured route.
// The Priority field indicates how soon the route works relative to other configured routes.
// Routes of equal priority are consulted in chronological order.
// The Description field provides a human-readable description for the route.
// Mailgun ignores this field except to provide the description when viewing the Mailgun web control panel.
// The Expression field lets you specify a pattern to match incoming messages against.
// The Actions field contains strings specifying what to do
// with any message which matches the provided expression.
// The CreatedAt field provides a time-stamp for when the route came into existence.
// Finally, the ID field provides a unique identifier for this route.
//
// When creating a new route, the SDK only uses a subset of the fields of this structure.
// In particular, CreatedAt and ID are meaningless in this context, and will be ignored.
// Only Priority, Description, Expression, and Actions need be provided.
type Route struct {
	Priority int `json:"priority,omitempty"`
	Description string `json:"description,omitempty"`
	Expression string `json:"expression,omitempty"`
	Actions []string `json:"actions,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`
	ID string `json:"id,omitempty"`
}

// GetRoutes returns the complete set of routes configured for your domain.
// You use routes to configure how to handle returned messages, or
// messages sent to a specfic address on your domain.
// See the Mailgun documentation for more information.
func (mg *mailgunImpl) GetRoutes(limit, skip int) (int, []Route, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(routesEndpoint))
	if limit != DefaultLimit {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		r.AddParameter("skip", strconv.Itoa(skip))
	}
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())

	var envelope struct{
		TotalCount int `json:"total_count"`
		Items []Route `json:"items"`
	}
	err := r.GetResponseFromJSON(&envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

// CreateRoute installs a new route for your domain.
// The route structure you provide serves as a template, and
// only a subset of the fields influence the operation.
// See the Route structure definition for more details.
func (mg *mailgunImpl) CreateRoute(prototype Route) (Route, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(routesEndpoint))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	p.AddValue("priority", strconv.Itoa(prototype.Priority))
	p.AddValue("description", prototype.Description)
	p.AddValue("expression", prototype.Expression)
	for _, action := range prototype.Actions {
		p.AddValue("action", action)
	}
	var envelope struct {
		Message string
		Route
	}
	_, err := r.MakePostRequest(p)
	return envelope.Route, err
}
