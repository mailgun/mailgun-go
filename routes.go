package mailgun

import (
	"context"
	"strconv"
)

// A Route structure contains information on a configured or to-be-configured route.
// When creating a new route, the SDK only uses a subset of the fields of this structure.
// In particular, CreatedAt and ID are meaningless in this context, and will be ignored.
// Only Priority, Description, Expression, and Actions need be provided.
type Route struct {
	// The Priority field indicates how soon the route works relative to other configured routes.
	// Routes of equal priority are consulted in chronological order.
	Priority int `json:"priority,omitempty"`
	// The Description field provides a human-readable description for the route.
	// Mailgun ignores this field except to provide the description when viewing the Mailgun web control panel.
	Description string `json:"description,omitempty"`
	// The Expression field lets you specify a pattern to match incoming messages against.
	Expression string `json:"expression,omitempty"`
	// The Actions field contains strings specifying what to do
	// with any message which matches the provided expression.
	Actions []string `json:"actions,omitempty"`

	// The CreatedAt field provides a time-stamp for when the route came into existence.
	CreatedAt RFC2822Time `json:"created_at,omitempty"`
	// ID field provides a unique identifier for this route.
	Id string `json:"id,omitempty"`
}

type routesListResponse struct {
	// is -1 if Next() or First() have not been called
	TotalCount int     `json:"total_count"`
	Items      []Route `json:"items"`
}

type createRouteResp struct {
	Message string `json:"message"`
	Route   `json:"route"`
}

// ListRoutes allows you to iterate through a list of routes returned by the API
func (mg *MailgunImpl) ListRoutes(opts *ListOptions) *RoutesIterator {
	var limit int
	if opts != nil {
		limit = opts.Limit
	}

	if limit == 0 {
		limit = 100
	}

	return &RoutesIterator{
		mg:                 mg,
		url:                generatePublicApiUrl(mg, routesEndpoint),
		routesListResponse: routesListResponse{TotalCount: -1},
		limit:              limit,
	}
}

type RoutesIterator struct {
	routesListResponse

	limit  int
	mg     Mailgun
	offset int
	url    string
	err    error
}

// If an error occurred during iteration `Err()` will return non nil
func (ri *RoutesIterator) Err() error {
	return ri.err
}

// Offset returns the current offset of the iterator
func (ri *RoutesIterator) Offset() int {
	return ri.offset
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ri *RoutesIterator) Next(ctx context.Context, items *[]Route) bool {
	if ri.err != nil {
		return false
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}

	cpy := make([]Route, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	if len(ri.Items) == 0 {
		return false
	}
	ri.offset = ri.offset + len(ri.Items)
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ri *RoutesIterator) First(ctx context.Context, items *[]Route) bool {
	if ri.err != nil {
		return false
	}
	ri.err = ri.fetch(ctx, 0, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Route, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	ri.offset = len(ri.Items)
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ri *RoutesIterator) Last(ctx context.Context, items *[]Route) bool {
	if ri.err != nil {
		return false
	}

	if ri.TotalCount == -1 {
		return false
	}

	ri.offset = ri.TotalCount - ri.limit
	if ri.offset < 0 {
		ri.offset = 0
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Route, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ri *RoutesIterator) Previous(ctx context.Context, items *[]Route) bool {
	if ri.err != nil {
		return false
	}

	if ri.TotalCount == -1 {
		return false
	}

	ri.offset = ri.offset - (ri.limit * 2)
	if ri.offset < 0 {
		ri.offset = 0
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Route, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	if len(ri.Items) == 0 {
		return false
	}
	return true
}

func (ri *RoutesIterator) fetch(ctx context.Context, skip, limit int) error {
	r := newHTTPRequest(ri.url)
	r.setBasicAuth(basicAuthUser, ri.mg.APIKey())
	r.setClient(ri.mg.Client())

	if skip != 0 {
		r.addParameter("skip", strconv.Itoa(skip))
	}
	if limit != 0 {
		r.addParameter("limit", strconv.Itoa(limit))
	}

	return getResponseFromJSON(ctx, r, &ri.routesListResponse)
}

// CreateRoute installs a new route for your domain.
// The route structure you provide serves as a template, and
// only a subset of the fields influence the operation.
// See the Route structure definition for more details.
func (mg *MailgunImpl) CreateRoute(ctx context.Context, prototype Route) (_ignored Route, err error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, routesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("priority", strconv.Itoa(prototype.Priority))
	p.addValue("description", prototype.Description)
	p.addValue("expression", prototype.Expression)
	for _, action := range prototype.Actions {
		p.addValue("action", action)
	}
	var resp createRouteResp
	if err = postResponseFromJSON(ctx, r, p, &resp); err != nil {
		return _ignored, err
	}
	return resp.Route, err
}

// DeleteRoute removes the specified route from your domain's configuration.
// To avoid ambiguity, Mailgun identifies the route by unique ID.
// See the Route structure definition and the Mailgun API documentation for more details.
func (mg *MailgunImpl) DeleteRoute(ctx context.Context, id string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, routesEndpoint) + "/" + id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// GetRoute retrieves the complete route definition associated with the unique route ID.
func (mg *MailgunImpl) GetRoute(ctx context.Context, id string) (Route, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, routesEndpoint) + "/" + id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		Message string `json:"message"`
		*Route  `json:"route"`
	}
	err := getResponseFromJSON(ctx, r, &envelope)
	if err != nil {
		return Route{}, err
	}
	return *envelope.Route, err

}

// UpdateRoute provides an "in-place" update of the specified route.
// Only those route fields which are non-zero or non-empty are updated.
// All other fields remain as-is.
func (mg *MailgunImpl) UpdateRoute(ctx context.Context, id string, route Route) (Route, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, routesEndpoint) + "/" + id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	if route.Priority != 0 {
		p.addValue("priority", strconv.Itoa(route.Priority))
	}
	if route.Description != "" {
		p.addValue("description", route.Description)
	}
	if route.Expression != "" {
		p.addValue("expression", route.Expression)
	}
	if route.Actions != nil {
		for _, action := range route.Actions {
			p.addValue("action", action)
		}
	}
	// For some reason, this API function just returns a bare Route on success.
	// Unsure why this is the case; it seems like it ought to be a bug.
	var envelope Route
	err := putResponseFromJSON(ctx, r, p, &envelope)
	return envelope, err
}
