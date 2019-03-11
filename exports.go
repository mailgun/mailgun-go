package mailgun

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type ExportList struct {
	Items []Export `json:"items"`
}

type Export struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

// Create an export based on the URL given
func (mg *MailgunImpl) CreateExport(ctx context.Context, url string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, exportsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("url", url)
	_, err := makePostRequest(ctx, r, payload)
	return err
}

// List all exports created within the past 24 hours
func (mg *MailgunImpl) ListExports(ctx context.Context, url string) ([]Export, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, exportsEndpoint))
	r.setClient(mg.Client())
	if url != "" {
		r.addParameter("url", url)
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp ExportList
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}

	var result []Export
	for _, item := range resp.Items {
		result = append(result, Export(item))
	}
	return result, nil
}

// GetExport gets an export by id
func (mg *MailgunImpl) GetExport(ctx context.Context, id string) (Export, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, exportsEndpoint) + "/" + id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp Export
	err := getResponseFromJSON(ctx, r, &resp)
	return resp, err
}

// Download an export by ID. This will respond with a '302 Moved'
// with the Location header of temporary S3 URL if it is available.
func (mg *MailgunImpl) GetExportLink(ctx context.Context, id string) (string, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, exportsEndpoint) + "/" + id + "/download_url")
	c := mg.Client()

	// Ensure the client doesn't attempt to retry
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("redirect")
	}

	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	r.addHeader("User-Agent", MailgunGoUserAgent)

	req, err := r.NewRequest(ctx, "GET", nil)
	if err != nil {
		return "", err
	}
	if Debug {
		fmt.Println(r.curlString(req, nil))
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusFound {
			url, err := resp.Location()
			if err != nil {
				return "", fmt.Errorf("while parsing 302 redirect url: %s", err)
			}
			return url.String(), nil
		}
		return "", err
	}
	return "", fmt.Errorf("expected a 302 response, API returned a '%d' instead", resp.StatusCode)
}
