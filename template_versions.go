package mailgun

import (
	"context"
	"strconv"
)

type TemplateVersion struct {
	Id        string         `json:"id"`
	Template  string         `json:"template,omitempty"`
	Engine    TemplateEngine `json:"engine"`
	CreatedAt string         `json:"createdAt"`
	//CreatedAt RFC2822Time    `json:"createdAt"`
	Comment string `json:"comment"`
	Active  bool   `json:"active"`
}

type templateVersionListResp struct {
	Item struct {
		Template
		Versions []TemplateVersion `json:"versions,omitempty"`
	} `json:"item"`
	Paging Paging `json:"paging"`
}

// Add a template version to a template
func (mg *MailgunImpl) AddTemplateVersion(ctx context.Context, templateId string, version *TemplateVersion) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateId + "/versions")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("template", version.Template)

	if version.Engine != "" {
		payload.addValue("engine", string(version.Engine))
	}
	if version.Comment != "" {
		payload.addValue("comment", version.Comment)
	}
	if version.Active {
		payload.addValue("active", boolToString(version.Active))
	}

	var resp templateResp
	if err := postResponseFromJSON(ctx, r, payload, &resp); err != nil {
		return err
	}
	*version = resp.Item.Version
	return nil
}

// Get a specific version of a template
func (mg *MailgunImpl) GetTemplateVersion(ctx context.Context, templateId, versionId string) (TemplateVersion, error) {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateId + "/versions/" + versionId)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp templateResp
	err := getResponseFromJSON(ctx, r, &resp)
	if err != nil {
		return TemplateVersion{}, err
	}
	return resp.Item.Version, nil
}

// Update the comment and mark a version of a template active
func (mg *MailgunImpl) UpdateTemplateVersion(ctx context.Context, templateId string, version *TemplateVersion) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateId + "/versions/" + version.Id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()

	if version.Comment != "" {
		p.addValue("comment", version.Comment)
	}
	if version.Active {
		p.addValue("active", boolToString(version.Active))
	}

	var resp templateResp
	err := putResponseFromJSON(ctx, r, p, &resp)
	if err != nil {
		return err
	}
	*version = resp.Item.Version
	return nil
}

// Delete a specific version of a template
func (mg *MailgunImpl) DeleteTemplateVersion(ctx context.Context, templateId, versionId string) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateId + "/versions/" + versionId)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

type TemplateVersionsIterator struct {
	templateVersionListResp
	mg  Mailgun
	err error
}

// List all the versions of a specific template
func (mg *MailgunImpl) ListTemplateVersions(templateId string, opts *ListOptions) *TemplateVersionsIterator {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateId + "/versions")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &TemplateVersionsIterator{
		mg:                      mg,
		templateVersionListResp: templateVersionListResp{Paging: Paging{Next: url, First: url}},
		err:                     err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (li *TemplateVersionsIterator) Err() error {
	return li.err
}

// Retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (li *TemplateVersionsIterator) Next(ctx context.Context, items *[]TemplateVersion) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Next)
	if li.err != nil {
		return false
	}
	cpy := make([]TemplateVersion, len(li.Item.Versions))
	copy(cpy, li.Item.Versions)
	*items = cpy
	if len(li.Item.Versions) == 0 {
		return false
	}
	return true
}

// Retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (li *TemplateVersionsIterator) First(ctx context.Context, items *[]TemplateVersion) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.First)
	if li.err != nil {
		return false
	}
	cpy := make([]TemplateVersion, len(li.Item.Versions))
	copy(cpy, li.Item.Versions)
	*items = cpy
	return true
}

// Retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (li *TemplateVersionsIterator) Last(ctx context.Context, items *[]TemplateVersion) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Last)
	if li.err != nil {
		return false
	}
	cpy := make([]TemplateVersion, len(li.Item.Versions))
	copy(cpy, li.Item.Versions)
	*items = cpy
	return true
}

// Retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (li *TemplateVersionsIterator) Previous(ctx context.Context, items *[]TemplateVersion) bool {
	if li.err != nil {
		return false
	}
	if li.Paging.Previous == "" {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Previous)
	if li.err != nil {
		return false
	}
	cpy := make([]TemplateVersion, len(li.Item.Versions))
	copy(cpy, li.Item.Versions)
	*items = cpy
	if len(li.Item.Versions) == 0 {
		return false
	}
	return true
}

func (li *TemplateVersionsIterator) fetch(ctx context.Context, url string) error {
	r := newHTTPRequest(url)
	r.setClient(li.mg.Client())
	r.setBasicAuth(basicAuthUser, li.mg.APIKey())

	return getResponseFromJSON(ctx, r, &li.templateVersionListResp)
}
