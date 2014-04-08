package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"fmt"
)

const (
	ReadOnly = "readonly"
	Members = "members"
	Everyone = "everyone"
)

type List struct {
	Address      string
	Name         string
	Description  string
	AccessLevel  string
	CreatedAt    string
	MembersCount int
}


func (mg *mailgunImpl) GetLists(limit, skip int, filter string) (int, []List, error) {
	r := simplehttp.NewHTTPRequest(generateApiUrl(mg, listsEndpoint))
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
	response, err := r.MakeRequest("GET", p)
	if err != nil {
		return -1, nil, err
	}
	fmt.Printf("@@ CODE(%d) DATA(%s)\n", response.Code, string(response.Data))
	return -1, nil, fmt.Errorf("Not finished")
}

func (mg *mailgunImpl) CreateList(prototype List) (List, error) {
	r := simplehttp.NewHTTPRequest(generateApiUrl(mg, listsEndpoint))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	p.AddValue("id", prototype.Address)
	p.AddValue("url", prototype.Name)
	_, err := r.MakePostRequest(p)
	return List{}, err
}

func (mg *mailgunImpl) DeleteList(addr string) error {
	r := simplehttp.NewHTTPRequest(generateApiUrl(mg, listsEndpoint) + "/ttt")
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	_, err := r.MakeDeleteRequest()
	return err
}

func (mg *mailgunImpl) GetListByAddress(addr string) (List, error) {
	r := simplehttp.NewHTTPRequest(generateApiUrl(mg, listsEndpoint) + "/ttt")
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	var envelope struct {
		List struct{
			Url *string `json:"url"`
		} `json:"List"`
	}
	err := r.GetResponseFromJSON(&envelope)
	return List{}, err
}

func (mg *mailgunImpl) UpdateList(addr string, prototype List) error {
	r := simplehttp.NewHTTPRequest(generateApiUrl(mg, listsEndpoint) + "/ttt")
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	p.AddValue("url", addr)
	_, err := r.MakePutRequest(p)
	return err
}