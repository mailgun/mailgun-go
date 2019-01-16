package mailgun

import (
	"context"
	"errors"
	"strconv"
)

type TemplateEngine string

// Used by CreateTemplate() and AddTemplateVersion() to specify the template engine
const (
	TemplateEngineMustache   = TemplateEngine("mustache")
	TemplateEngineHandlebars = TemplateEngine("handlebars")
	TemplateEngineGo         = TemplateEngine("go")
)

type Template struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   RFC2822Time `json:"createdAt"`
	Version   TemplateVersion `json:"version,omitempty"`
}

type templateResp struct {
	Item    Template `json:"item"`
	Message string   `json:"message"`
}

type templateListResp struct {
	Items  []Template `json:"items"`
	Paging Paging     `json:"paging"`
}

// Create a new template which can be used to attach template versions to
func (mg *MailgunImpl) CreateTemplate(ctx context.Context, template *Template) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()

	if template.Name != "" {
		payload.addValue("name", template.Name)
	}
	if template.Description != "" {
		payload.addValue("description", template.Description)
	}

	if template.Version.Engine != "" {
		payload.addValue("engine", string(template.Version.Engine))
	}
	if template.Version.Template != "" {
		payload.addValue("template", template.Version.Template)
	}
	if template.Version.Comment != "" {
		payload.addValue("comment", template.Version.Comment)
	}

	var resp templateResp
	if err := postResponseFromJSON(ctx, r, payload, &resp); err != nil {
		return err
	}
	*template = resp.Item
	return nil
}

// Get a template given the template id
func (mg *MailgunImpl) GetTemplate(ctx context.Context, id string) (Template, error) {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addParameter("active", "yes")

	var resp templateResp
	err := getResponseFromJSON(ctx, r, &resp)
	if err != nil {
		return Template{}, err
	}
	return resp.Item, nil
}

// Update the name and description of a template
func (mg *MailgunImpl) UpdateTemplate(ctx context.Context, template *Template) error {
	if template.Id == "" {
		return errors.New("UpdateTemplate() Template.Id cannot be empty")
	}

	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + template.Id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()

	if template.Name != "" {
		p.addValue("name", template.Name)
	}
	if template.Description != "" {
		p.addValue("description", template.Description)
	}

	var resp templateResp
	err := putResponseFromJSON(ctx, r, p, &resp)
	if err != nil {
		return err
	}
	*template = resp.Item
	return nil
}

// Delete a template given a template id
func (mg *MailgunImpl) DeleteTemplate(ctx context.Context, id string) error {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint) + "/" + id)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

type TemplatesIterator struct {
	templateListResp
	mg  Mailgun
	err error
}

// List all available templates
func (mg *MailgunImpl) ListTemplates(opts *ListOptions) *TemplatesIterator {
	r := newHTTPRequest(generateApiUrl(mg, templatesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &TemplatesIterator{
		mg:               mg,
		templateListResp: templateListResp{Paging: Paging{Next: url, First: url}},
		err:              err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (ti *TemplatesIterator) Err() error {
	return ti.err
}

// Retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ti *TemplatesIterator) Next(ctx context.Context, items *[]Template) bool {
	if ti.err != nil {
		return false
	}
	ti.err = ti.fetch(ctx, ti.Paging.Next)
	if ti.err != nil {
		return false
	}
	cpy := make([]Template, len(ti.Items))
	copy(cpy, ti.Items)
	*items = cpy
	if len(ti.Items) == 0 {
		return false
	}
	return true
}

// Retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ti *TemplatesIterator) First(ctx context.Context, items *[]Template) bool {
	if ti.err != nil {
		return false
	}
	ti.err = ti.fetch(ctx, ti.Paging.First)
	if ti.err != nil {
		return false
	}
	cpy := make([]Template, len(ti.Items))
	copy(cpy, ti.Items)
	*items = cpy
	return true
}

// Retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ti *TemplatesIterator) Last(ctx context.Context, items *[]Template) bool {
	if ti.err != nil {
		return false
	}
	ti.err = ti.fetch(ctx, ti.Paging.Last)
	if ti.err != nil {
		return false
	}
	cpy := make([]Template, len(ti.Items))
	copy(cpy, ti.Items)
	*items = cpy
	return true
}

// Retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ti *TemplatesIterator) Previous(ctx context.Context, items *[]Template) bool {
	if ti.err != nil {
		return false
	}
	if ti.Paging.Previous == "" {
		return false
	}
	ti.err = ti.fetch(ctx, ti.Paging.Previous)
	if ti.err != nil {
		return false
	}
	cpy := make([]Template, len(ti.Items))
	copy(cpy, ti.Items)
	*items = cpy
	if len(ti.Items) == 0 {
		return false
	}
	return true
}

func (ti *TemplatesIterator) fetch(ctx context.Context, url string) error {
	r := newHTTPRequest(url)
	r.setClient(ti.mg.Client())
	r.setBasicAuth(basicAuthUser, ti.mg.APIKey())

	return getResponseFromJSON(ctx, r, &ti.templateListResp)
}
