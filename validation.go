package mailgun

import (
	"context"
	"fmt"
)

// ValidateEmailResponse records basic facts about a validated e-mail address.
// See the ValidateEmail method and example for more details.
type ValidateEmailResponse struct {
	// Echoes the address provided.
	Address string `json:"address"`

	// Indicates whether Mailgun thinks the address is from a known
	// disposable mailbox provider.
	IsDisposableAddress bool `json:"is_disposable_address"`

	// Indicates whether Mailgun thinks the address is an email distribution list.
	IsRoleAddress bool `json:"is_role_address"`

	// A list of potential reasons why a specific validation may be unsuccessful. (Available in the v4 response)
	Reason []string `json:"reason"`

	// Result
	Result string `json:"result"`

	// Risk assessment for the provided email.
	Risk string `json:"risk"`

	LastSeen int64 `json:"last_seen,omitempty"`

	// Provides a simple recommendation in case the address is invalid or
	// Mailgun thinks you might have a typo. May be empty, in which case
	// Mailgun has no recommendation to give.
	DidYouMean string `json:"did_you_mean,omitempty"`

	// Engagement results are a macro-level view that explain an email recipient’s propensity to engage.
	// https://documentation.mailgun.com/docs/inboxready/mailgun-validate/validate_engagement/
	Engagement *EngagementData `json:"engagement,omitempty"`

	RootAddress string `json:"root_address,omitempty"`
}

type EngagementData struct {
	Engaging bool   `json:"engaging"`
	IsBot    bool   `json:"is_bot"`
	Behavior string `json:"behavior,omitempty"`
}

// ValidateEmail performs various checks on the email address provided to ensure it's correctly formatted.
// It may also be used to break an email address into its sub-components.
// https://documentation.mailgun.com/docs/inboxready/mailgun-validate/single-valid-ir/
func (mg *MailgunImpl) ValidateEmail(ctx context.Context, email string, mailBoxVerify bool) (ValidateEmailResponse, error) {
	r := newHTTPRequest(fmt.Sprintf("%s/v4/address/validate", mg.APIBase()))
	r.setClient(mg.HTTPClient())
	r.addParameter("address", email)
	if mailBoxVerify {
		r.addParameter("mailbox_verification", "true")
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var res ValidateEmailResponse
	err := getResponseFromJSON(ctx, r, &res)
	if err != nil {
		return ValidateEmailResponse{}, err
	}

	return res, nil
}
