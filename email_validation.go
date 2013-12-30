package mailgun

import (
	"github.com/mbanzon/simplehttp"
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
	r := simplehttp.NewGetRequest(generatePublicApiUrl(addressValidateEndpoint))
	r.AddParameter("address", email)
	r.SetBasicAuth(basicAuthUser, m.PublicApiKey())

	var response EmailVerification
	err := r.MakeJSONRequest(&response)
	if err != nil {
		return EmailVerification{}, err
	}

	return response, nil
}

func (m *mailgunImpl) ParseAddresses(addresses ...string) ([]string, []string, error) {
	r := simplehttp.NewGetRequest(generatePublicApiUrl(addressParseEndpoint))
	r.AddParameter("addresses", strings.Join(addresses, ","))
	r.SetBasicAuth(basicAuthUser, m.PublicApiKey())

	var response AddressParseResult
	err := r.MakeJSONRequest(&response)
	if err != nil {
		return nil, nil, err
	}

	return response.Parsed, response.Unparseable, nil
}
