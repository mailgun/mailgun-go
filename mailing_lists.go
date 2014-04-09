package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"encoding/json"
//	"fmt"
)

const (
	ReadOnly = "readonly"
	Members = "members"
	Everyone = "everyone"
)

type Subscription string

var (
	All *bool = nil
	Subscribed *bool = &yes
	Unsubscribed *bool = &no
)

var (
	yes bool = true
	no bool = false
)

type List struct {
	Address      string `json:"address",omitempty"`
	Name         string `json:"name",omitempty"`
	Description  string `json:"description",omitempty"`
	AccessLevel  string `json:"access_level",omitempty"`
	CreatedAt    string `json:"created_at",omitempty"`
	MembersCount int    `json:"members_count",omitempty"`
}

type Subscriber struct {
	Address    string `json:"address,omitempty"`
	Name       string `json:"name,omitempty"`
	Subscribed *bool `json:"subscribed,omitempty"`
	Vars       map[string]interface{} `json:"vars,omitempty"`
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


func (mg *mailgunImpl) GetSubscribers(limit, skip int, s *bool, addr string) (int, []Subscriber, error) {
	r := simplehttp.NewHTTPRequest(generateSubscriberApiUrl(listsEndpoint, addr))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	if limit != DefaultLimit {
		p.AddValue("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		p.AddValue("skip", strconv.Itoa(skip))
	}
	if s != nil {
		p.AddValue("subscribed", yesNo(*s))
	}
	var envelope struct {
		TotalCount int `json:"total_count"`
		Items []Subscriber `json:"items"`
	}
	response, err := r.MakeRequest("GET", p)
	if err != nil {
		return -1, nil, err
	}
	err = response.ParseFromJSON(&envelope)
	return envelope.TotalCount, envelope.Items, err
}

func (mg *mailgunImpl) GetSubscriberByAddress(s, l string) (Subscriber, error) {
	r := simplehttp.NewHTTPRequest(generateSubscriberApiUrl(listsEndpoint, l) + "/" + s)
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	response, err := r.MakeGetRequest()
	if err != nil {
		return Subscriber{}, err
	}
	var envelope struct {
		Member Subscriber `json:"member"`
	}
	err = response.ParseFromJSON(&envelope)
	return envelope.Member, err
}

func (mg *mailgunImpl) CreateSubscriber(merge bool, addr string, prototype Subscriber) error {
	vs, err := json.Marshal(prototype.Vars)
	if err != nil {
		return err
	}

	r := simplehttp.NewHTTPRequest(generateSubscriberApiUrl(listsEndpoint, addr))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewFormDataPayload()
	p.AddValue("upsert", yesNo(merge))
	p.AddValue("address", prototype.Address)
	p.AddValue("name", prototype.Name)
	p.AddValue("vars", string(vs))
	if prototype.Subscribed != nil {
		p.AddValue("subscribed", yesNo(*prototype.Subscribed))
	}
	_, err = r.MakePostRequest(p)
	return err
}

func (mg *mailgunImpl) UpdateSubscriber(s, l string, prototype Subscriber) (Subscriber, error) {
	r := simplehttp.NewHTTPRequest(generateSubscriberApiUrl(listsEndpoint, l) + "/" + s)
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewFormDataPayload()
	if prototype.Address != "" {
		p.AddValue("address", prototype.Address)
	}
	if prototype.Name != "" {
		p.AddValue("name", prototype.Name)
	}
	if prototype.Vars != nil {
		vs, err := json.Marshal(prototype.Vars)
		if err != nil {
			return Subscriber{}, err
		}
		p.AddValue("vars", string(vs))
	}
	if prototype.Subscribed != nil {
		p.AddValue("subscribed", yesNo(*prototype.Subscribed))
	}
	response, err := r.MakePutRequest(p)
	if err != nil {
		return Subscriber{}, err
	}
	var envelope struct {
		Member Subscriber `json:"member"`
	}
	err = response.ParseFromJSON(&envelope)
	return envelope.Member, err
}
