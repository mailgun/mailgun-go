package mailgun

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type EmailVarificationParts struct {
	LocalPart   string `json:"local_part"`
	Domain      string `json:"domain"`
	DisplayName string `json:"display_name"`
}

type EmailVerification struct {
	IsValid    bool                   `json:"is_valid"`
	Parts      EmailVarificationParts `json:"parts"`
	Address    string                 `json:"address"`
	DidYouMean string                 `json:"did_you_mean"`
}

type AddressParseResult struct {
	Parsed      []string `json:"parsed"`
	Unparseable []string `json:"unparseable"`
}

func (m *mailgunImpl) ValidateEmail(email string) (EmailVerification, error) {
	u, err := url.Parse(generatePublicApiUrl(addressValidateEndpoint))
	if err != nil {
		return EmailVerification{}, err
	}
	q := u.Query()
	q.Set("address", email)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return EmailVerification{}, err
	}
	req.SetBasicAuth(basicAuthUser, m.PublicApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return EmailVerification{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return EmailVerification{}, errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EmailVerification{}, err
	}

	var response EmailVerification
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		return EmailVerification{}, err2
	}

	return response, nil
}

func (m *mailgunImpl) ParseAddresses(addresses ...string) ([]string, []string, error) {
	u, err := url.Parse(generatePublicApiUrl(addressParseEndpoint))
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	q.Set("addresses", strings.Join(addresses, ","))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, nil, err
	}
	req.SetBasicAuth(basicAuthUser, m.PublicApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var response AddressParseResult
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		return nil, nil, err2
	}

	return response.Parsed, response.Unparseable, nil
}
