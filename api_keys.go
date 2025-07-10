package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type ListAPIKeysOptions struct {
	DomainName string
	Kind       string
}

func (mg *Client) ListAPIKeys(ctx context.Context, opts *ListAPIKeysOptions) ([]mtypes.APIKey, error) {
	r := newHTTPRequest(generateApiUrl(mg, mtypes.APIKeysVersion, mtypes.APIKeysEndpoint))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	if opts != nil {
		if opts.DomainName != "" {
			r.addParameter("domain_name", opts.DomainName)
		}

		if opts.Kind != "" {
			r.addParameter("kind", opts.Kind)
		}
	}

	var resp mtypes.APIKeyList
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}

	var result []mtypes.APIKey
	for _, item := range resp.Items {
		result = append(result, item)
	}
	return result, nil
}
