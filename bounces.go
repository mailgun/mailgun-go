package mailgun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BounceItem struct {
	CreatedAt string `json:"created_at"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	Error     string `json:"error"`
}

type Bounces struct {
	TotalCount int          `json:"total_count"`
	Items      []BounceItem `json:"items"`
}

type singleBounce struct {
	Bounce BounceItem `json:"bounce"`
}

func (i BounceItem) GetCreatedAt() (t time.Time, err error) {
	t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", i.CreatedAt)
	return
}

func (m *mailgunImpl) GetBounces(limit, skip int) (Bounces, error) {
	u, err := url.Parse(generateApiUrl(m, bouncesEndpoint))
	if err != nil {
		return Bounces{}, err
	}

	q := u.Query()
	if limit != -1 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		q.Set("skip", strconv.Itoa(skip))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Bounces{}, err
	}
	req.SetBasicAuth(basicAuthUser, m.ApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Bounces{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Bounces{}, errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Bounces{}, err
	}

	var response Bounces
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		return Bounces{}, err2
	}

	return response, nil
}

func (m *mailgunImpl) GetSingleBounce(address string) (BounceItem, error) {
	u, err := url.Parse(generateApiUrl(m, bouncesEndpoint) + "/" + address)
	if err != nil {
		return BounceItem{}, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return BounceItem{}, err
	}
	req.SetBasicAuth(basicAuthUser, m.ApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BounceItem{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return BounceItem{}, errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return BounceItem{}, err
	}

	var response singleBounce
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		return BounceItem{}, err2
	}

	return response.Bounce, nil
}

func (m *mailgunImpl) AddBounce(address, code, error string) error {
	data := url.Values{}
	data.Add("address", address)
	if code != "" {
		data.Add("code", code)
	}
	if error != "" {
		data.Add("error", error)
	}

	u, err := url.Parse(generateApiUrl(m, bouncesEndpoint))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(basicAuthUser, m.ApiKey())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	return nil
}

func (m *mailgunImpl) DeleteBounce(address string) error {
	u, err := url.Parse(generateApiUrl(m, bouncesEndpoint) + "/" + address)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(basicAuthUser, m.ApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	return nil
}
