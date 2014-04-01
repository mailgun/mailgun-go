package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"fmt"
	"strconv"
)

type Route struct {
	Priority int `json:"priority,omitempty"`
	Description string `json:"description,omitempty"`
	Expression string `json:"expression,omitempty"`
	Actions []string `json:"actions,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`
	ID string `json:"id,omitempty"`
}

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

func (mg *mailgunImpl) CreateRoute(prototype Route) (Route, error) {
	return Route{}, fmt.Errorf("CreateRoute: not implemented")
}
