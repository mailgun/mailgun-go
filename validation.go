package mailgun

import (
	"context"
	"fmt"
	"strings"

	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v4/mtypes"
)

// ValidateEmail performs various checks on the email address provided to ensure it's correctly formatted.
// It may also be used to break an email address into its sub-components.
// https://documentation.mailgun.com/docs/inboxready/mailgun-validate/single-valid-ir/
func (mg *MailgunImpl) ValidateEmail(ctx context.Context, email string, mailBoxVerify bool) (mtypes.ValidateEmailResponse, error) {
	if !strings.HasSuffix(mg.APIBase(), "/v4") {
		return mtypes.ValidateEmailResponse{}, errors.New("ValidateEmail: only v4 is supported")
	}

	r := newHTTPRequest(fmt.Sprintf("%s/address/validate", mg.APIBase()))
	r.setClient(mg.Client())
	r.addParameter("address", email)
	if mailBoxVerify {
		r.addParameter("mailbox_verification", "true")
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var res mtypes.ValidateEmailResponse
	err := getResponseFromJSON(ctx, r, &res)
	if err != nil {
		return mtypes.ValidateEmailResponse{}, err
	}

	return res, nil
}
