package mailgun

import (
	"context"
	"strconv"
)

type TemplateVersion struct {
	Tag       string         `json:"tag"`
	Template  string         `json:"template,omitempty"`
	Engine    TemplateEngine `json:"engine"`
	CreatedAt RFC2822Time    `json:"createdAt"`
	Comment   string         `json:"comment"`
	Active    bool           `json:"active"`
}

type templateVersionListResp struct {
	Template struct {
		Template
		Versions []TemplateVersion `json:"versions,omitempty"`
	} `json:"template"`
	Paging Paging `json:"paging"`
}

// AddTemplateVersion adds a template version to a template
func (mg *MailgunImpl) AddTemplateVersion(ctx context.Context, templateName string, version *TemplateVersion) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateName + "/versions")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("template", version.Template)

	if version.Tag != "" {
		payload.addValue("tag", string(version.Tag))
	}
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

// GetTemplateVersion gets a specific version of a template
func (mg *MailgunImpl) GetTemplateVersion(ctx context.Context, templateName, tag string) (TemplateVersion, error) {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateName + "/versions/" + tag)
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
func (mg *MailgunImpl) UpdateTemplateVersion(ctx context.Context, templateName string, version *TemplateVersion) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateName + "/versions/" + version.Tag)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()

	if version.Comment != "" {
		p.addValue("comment", version.Comment)
	}
	if version.Active {
		p.addValue("active", boolToString(version.Active))
	}
	if version.Template != "" {
		p.addValue("template", version.Template)
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
func (mg *MailgunImpl) DeleteTemplateVersion(ctx context.Context, templateName, tag string) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateName + "/versions/" + tag)
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
func (mg *MailgunImpl) ListTemplateVersions(templateName string, opts *ListOptions) *TemplateVersionsIterator {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + templateName + "/versions")
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

// Next retrieves the next page of items from the api. Returns false when there
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
	cpy := make([]TemplateVersion, len(li.Template.Versions))
	copy(cpy, li.Template.Versions)
	*items = cpy
	if len(li.Template.Versions) == 0 {
		return false
	}
	return true
}

// First retrieves the first page of items from the api. Returns false if there
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
	cpy := make([]TemplateVersion, len(li.Template.Versions))
	copy(cpy, li.Template.Versions)
	*items = cpy
	return true
}

// Last retrieves the last page of items from the api.
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
	cpy := make([]TemplateVersion, len(li.Template.Versions))
	copy(cpy, li.Template.Versions)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
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
	cpy := make([]TemplateVersion, len(li.Template.Versions))
	copy(cpy, li.Template.Versions)
	*items = cpy
	if len(li.Template.Versions) == 0 {
		return false
	}
	return true
}

func (li *TemplateVersionsIterator) fetch(ctx context.Context, url string) error {
	li.Template.Versions = nil
	r := newHTTPRequest(url)
	r.setClient(li.mg.Client())
	r.setBasicAuth(basicAuthUser, li.mg.APIKey())

	return getResponseFromJSON(ctx, r, &li.templateVersionListResp)
}
