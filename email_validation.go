package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strings"
)

// TODO(sfalvo): Document me.
type EmailVerificationParts struct {
	LocalPart   string `json:"local_part"`
	Domain      string `json:"domain"`
	DisplayName string `json:"display_name"`
}

// TODO(sfalvo): Document me.
type EmailVerification struct {
	IsValid    bool                   `json:"is_valid"`
	Parts      EmailVerificationParts `json:"parts"`
	Address    string                 `json:"address"`
	DidYouMean string                 `json:"did_you_mean"`
}

// TODO(sfalvo): Document me.
type AddressParseResult struct {
	Parsed      []string `json:"parsed"`
	Unparseable []string `json:"unparseable"`
}

// ValidateEmail performs various checks on the email address provided to ensure it's correctly formatted.
// It may also be used to break an email address into its sub-components.  (See example.)
// NOTE: Use of this function requires a proper public API key.  The private API key will not work.
func (m *MailgunImpl) ValidateEmail(email string) (EmailVerification, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(addressValidateEndpoint))
	r.AddParameter("address", email)
	r.SetBasicAuth(basicAuthUser, m.PublicApiKey())

	var response EmailVerification
	err := getResponseFromJSON(r, &response)
	if err != nil {
		return EmailVerification{}, err
	}

	return response, nil
}

// ParseAddresses takes a list of addresses and sorts them into valid and invalid address categories.
// NOTE: Use of this function requires a proper public API key.  The private API key will not work.
func (m *MailgunImpl) ParseAddresses(addresses ...string) ([]string, []string, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(addressParseEndpoint))
	r.AddParameter("addresses", strings.Join(addresses, ","))
	r.SetBasicAuth(basicAuthUser, m.PublicApiKey())

	var response AddressParseResult
	err := getResponseFromJSON(r, &response)
	if err != nil {
		return nil, nil, err
	}

	return response.Parsed, response.Unparseable, nil
}
