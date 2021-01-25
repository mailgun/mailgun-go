package mailgun

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// The EmailVerificationParts structure breaks out the basic elements of an email address.
// LocalPart includes everything up to the '@' in an e-mail address.
// Domain includes everything after the '@'.
// DisplayName is no longer used, and will appear as "".
type EmailVerificationParts struct {
	LocalPart   string `json:"local_part"`
	Domain      string `json:"domain"`
	DisplayName string `json:"display_name"`
}

// EmailVerification records basic facts about a validated e-mail address.
// See the ValidateEmail method and example for more details.
//
type EmailVerification struct {
	// Indicates whether an email address conforms to IETF RFC standards.
	IsValid bool `json:"is_valid"`
	// Indicates whether an email address is deliverable.
	MailboxVerification string `json:"mailbox_verification"`
	// Parts records the different subfields of the parsed email address
	Parts EmailVerificationParts `json:"parts"`
	// Echoes the address provided.
	Address string `json:"address"`
	// Provides a simple recommendation in case the address is invalid or
	// Mailgun thinks you might have a typo. May be empty, in which case
	// Mailgun has no recommendation to give.
	DidYouMean string `json:"did_you_mean"`
	// Indicates whether Mailgun thinks the address is from a known
	// disposable mailbox provider.
	IsDisposableAddress bool `json:"is_disposable_address"`
	// Indicates whether Mailgun thinks the address is an email distribution list.
	IsRoleAddress bool `json:"is_role_address"`
	// The reason why a specific validation may be unsuccessful. (Available in the V3 response)
	Reason string `json:"reason"`
	// A list of potential reasons why a specific validation may be unsuccessful. (Available in the v4 response)
	Reasons []string
	// Risk assessment for the provided email.
	Risk string `json:"risk"`
}

type v4EmailValidationResp struct {
	IsValid             bool                   `json:"is_valid"`
	MailboxVerification string                 `json:"mailbox_verification"`
	Parts               EmailVerificationParts `json:"parts"`
	Address             string                 `json:"address"`
	DidYouMean          string                 `json:"did_you_mean"`
	IsDisposableAddress bool                   `json:"is_disposable_address"`
	IsRoleAddress       bool                   `json:"is_role_address"`
	Reason              []string               `json:"reason"`
	Risk                string                 `json:"risk"`
}

type addressParseResult struct {
	Parsed      []string `json:"parsed"`
	Unparseable []string `json:"unparseable"`
}

type EmailValidator interface {
	ValidateEmail(ctx context.Context, email string, mailBoxVerify bool) (EmailVerification, error)
	ParseAddresses(ctx context.Context, addresses ...string) ([]string, []string, error)
}

type EmailValidatorImpl struct {
	client      *http.Client
	isPublicKey bool
	apiBase     string
	apiKey      string
}

// Creates a new validation instance.
// * If a public key is provided, uses the public validation endpoints
// * If a private key is provided, uses the private validation endpoints
func NewEmailValidator(apiKey string) *EmailValidatorImpl {
	isPublicKey := false

	// Did the user pass in a public key?
	if strings.HasPrefix(apiKey, "pubkey-") {
		isPublicKey = true
	}

	return &EmailValidatorImpl{
		client:      http.DefaultClient,
		isPublicKey: isPublicKey,
		apiBase:     APIBase,
		apiKey:      apiKey,
	}
}

// NewEmailValidatorFromEnv returns a new EmailValidator using environment variables
// If MG_PUBLIC_API_KEY is set, assume using the free validation subject to daily usage limits
// If only MG_API_KEY is set, assume using the /private validation routes with no daily usage limits
func NewEmailValidatorFromEnv() (*EmailValidatorImpl, error) {
	apiKey := os.Getenv("MG_PUBLIC_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MG_API_KEY")
		if apiKey == "" {
			return nil, errors.New(
				"environment variable MG_PUBLIC_API_KEY or MG_API_KEY required for email validation")
		}
	}
	v := NewEmailValidator(apiKey)
	url := os.Getenv("MG_URL")
	if url != "" {
		v.SetAPIBase(url)
	}
	return v, nil
}

// APIBase returns the API Base URL configured for this client.
func (m *EmailValidatorImpl) APIBase() string {
	return m.apiBase
}

// SetAPIBase updates the API Base URL for this client.
func (m *EmailValidatorImpl) SetAPIBase(address string) {
	m.apiBase = address
}

// SetClient updates the HTTP client for this client.
func (m *EmailValidatorImpl) SetClient(c *http.Client) {
	m.client = c
}

// Client returns the HTTP client configured for this client.
func (m *EmailValidatorImpl) Client() *http.Client {
	return m.client
}

// APIKey returns the API key used for validations
func (m *EmailValidatorImpl) APIKey() string {
	return m.apiKey
}

func (m *EmailValidatorImpl) getAddressURL(endpoint string) string {
	if m.isPublicKey {
		return fmt.Sprintf("%s/address/%s", m.APIBase(), endpoint)
	}
	return fmt.Sprintf("%s/address/private/%s", m.APIBase(), endpoint)
}

// ValidateEmail performs various checks on the email address provided to ensure it's correctly formatted.
// It may also be used to break an email address into its sub-components. If user has set the
func (m *EmailValidatorImpl) ValidateEmail(ctx context.Context, email string, mailBoxVerify bool) (EmailVerification, error) {
	if strings.HasSuffix(m.APIBase(), "/v4") {
		return m.validateV4(ctx, email, mailBoxVerify)
	}
	return m.validateV3(ctx, email, mailBoxVerify)
}

func (m *EmailValidatorImpl) validateV3(ctx context.Context, email string, mailBoxVerify bool) (EmailVerification, error) {
	r := newHTTPRequest(m.getAddressURL("validate"))
	r.setClient(m.Client())
	r.addParameter("address", email)
	if mailBoxVerify {
		r.addParameter("mailbox_verification", "true")
	}
	r.setBasicAuth(basicAuthUser, m.APIKey())

	var res EmailVerification
	err := getResponseFromJSON(ctx, r, &res)
	if err != nil {
		return EmailVerification{}, err
	}
	return res, nil
}

func (m *EmailValidatorImpl) validateV4(ctx context.Context, email string, mailBoxVerify bool) (EmailVerification, error) {
	r := newHTTPRequest(fmt.Sprintf("%s/address/validate", m.APIBase()))
	r.setClient(m.Client())
	r.addParameter("address", email)
	if mailBoxVerify {
		r.addParameter("mailbox_verification", "true")
	}
	r.setBasicAuth(basicAuthUser, m.APIKey())

	var res v4EmailValidationResp
	err := getResponseFromJSON(ctx, r, &res)
	if err != nil {
		return EmailVerification{}, err
	}
	return EmailVerification{
		IsValid:             res.IsValid,
		MailboxVerification: res.MailboxVerification,
		Parts:               res.Parts,
		Address:             res.Address,
		DidYouMean:          res.DidYouMean,
		IsDisposableAddress: res.IsDisposableAddress,
		IsRoleAddress:       res.IsRoleAddress,
		Reasons:             res.Reason,
		Risk:                res.Risk}, nil
}

// ParseAddresses takes a list of addresses and sorts them into valid and invalid address categories.
// NOTE: Use of this function requires a proper public API key.  The private API key will not work.
func (m *EmailValidatorImpl) ParseAddresses(ctx context.Context, addresses ...string) ([]string, []string, error) {
	r := newHTTPRequest(m.getAddressURL("parse"))
	r.setClient(m.Client())
	r.addParameter("addresses", strings.Join(addresses, ","))
	r.setBasicAuth(basicAuthUser, m.APIKey())

	var response addressParseResult
	err := getResponseFromJSON(ctx, r, &response)
	if err != nil {
		return nil, nil, err
	}

	return response.Parsed, response.Unparseable, nil
}
