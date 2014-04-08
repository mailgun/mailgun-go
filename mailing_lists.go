package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
)

const (
	ReadOnly = "readonly"
	Members = "members"
	Everyone = "everyone"
)

type List struct {
	Address      string `json:"address",omitempty"`
	Name         string `json:"name",omitempty"`
	Description  string `json:"description",omitempty"`
	AccessLevel  string `json:"access_level",omitempty"`
	CreatedAt    string `json:"created_at",omitempty"`
	MembersCount int    `json:"members_count",omitempty"`
}


func (mg *mailgunImpl) GetLists(limit, skip int, filter string) (int, []List, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(listsEndpoint))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	if limit != DefaultLimit {
		p.AddValue("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		p.AddValue("skip", strconv.Itoa(skip))
	}
	if filter != "" {
		p.AddValue("address", filter)
	}
	var envelope struct {
		Items []List `json:"items"`
		TotalCount int `json:"total_count"`
	}
	response, err := r.MakeRequest("GET", p)
	if err != nil {
		return -1, nil, err
	}
	err = response.ParseFromJSON(&envelope)
	return envelope.TotalCount, envelope.Items, err
}

func (mg *mailgunImpl) CreateList(prototype List) (List, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(listsEndpoint))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	if prototype.Address != "" {
		p.AddValue("address", prototype.Address)
	}
	if prototype.Name != "" {
		p.AddValue("name", prototype.Name)
	}
	if prototype.Description != "" {
		p.AddValue("description", prototype.Description)
	}
	if prototype.AccessLevel != "" {
		p.AddValue("access_level", prototype.AccessLevel)
	}
	response, err := r.MakePostRequest(p)
	if err != nil {
		return List{}, err
	}
	var l List
	err = response.ParseFromJSON(&l)
	return l, err
}

func (mg *mailgunImpl) DeleteList(addr string) error {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(listsEndpoint) + "/" + addr)
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	_, err := r.MakeDeleteRequest()
	return err
}

func (mg *mailgunImpl) GetListByAddress(addr string) (List, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(listsEndpoint) + "/" + addr)
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	response, err := r.MakeGetRequest()
	var envelope struct {
		List `json:"list"`
	}
	err = response.ParseFromJSON(&envelope)
	return envelope.List, err
}

func (mg *mailgunImpl) UpdateList(addr string, prototype List) (List, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(listsEndpoint) + "/" + addr)
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	if prototype.Address != "" {
		p.AddValue("address", prototype.Address)
	}
	if prototype.Name != "" {
		p.AddValue("name", prototype.Name)
	}
	if prototype.Description != "" {
		p.AddValue("description", prototype.Description)
	}
	if prototype.AccessLevel != "" {
		p.AddValue("access_level", prototype.AccessLevel)
	}
	var l List
	response, err := r.MakePutRequest(p)
	if err != nil {
		return l, err
	}
	err = response.ParseFromJSON(&l)
	return l, err
}
