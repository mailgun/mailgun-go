package mailgun

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mailgun/errors"
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
	// Result
	Result string `json:"result"`
	// Engagement results are a macro-level view that explain an email recipient’s propensity to engage.
	// https://documentation.mailgun.com/docs/inboxready/mailgun-validate/validate_engagement/
	Engagement *EngagementData `json:"engagement,omitempty"`
}

type EngagementData struct {
	Engaging bool   `json:"engaging"`
	Behavior string `json:"behavior,omitempty"`
	IsBot    bool   `json:"is_bot"`
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
	Result              string                 `json:"result"`
	Engagement          *EngagementData        `json:"engagement,omitempty"`
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
	client  *http.Client
	apiBase string
	apiKey  string
}

// NewEmailValidator creates a new validation instance.
// Use private key.
func NewEmailValidator(apiKey string) *EmailValidatorImpl {
	return &EmailValidatorImpl{
		// TODO(vtopc): Don’t use http.DefaultClient - https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
		client:  http.DefaultClient,
		apiBase: "https://api.mailgun.net/v4",
		apiKey:  apiKey,
	}
}

// NewEmailValidatorFromEnv returns a new EmailValidator using environment variables
//
// Set MG_API_KEY
func NewEmailValidatorFromEnv() (*EmailValidatorImpl, error) {
	apiKey := os.Getenv("MG_API_KEY")
	if apiKey == "" {
		return nil, errors.New(
			"environment variable MG_API_KEY required for email validation")
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

// ValidateEmail performs various checks on the email address provided to ensure it's correctly formatted.
// It may also be used to break an email address into its sub-components. If user has set the
// TODO(DE-1384): move to *MailgunImpl?
func (m *EmailValidatorImpl) ValidateEmail(ctx context.Context, email string, mailBoxVerify bool) (EmailVerification, error) {
	// TODO(DE-1383): remove check:
	if strings.HasSuffix(m.APIBase(), "/v4") {
		return m.validateV4(ctx, email, mailBoxVerify)
	}

	return EmailVerification{}, errors.New("ValidateEmail: only v4 is supported")
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
		Result:              res.Result,
		Risk:                res.Risk,
		Engagement:          res.Engagement,
	}, nil
}
