package mocks

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addTemplateRoutes(r chi.Router) {
	r.Get("/{domain}/templates", ms.listTemplates)
	r.Get("/{domain}/templates/{name}", ms.getTemplate)

	r.Post("/{domain}/templates", ms.createTemplate)
	r.Put("/{domain}/templates/{name}", ms.updateTemplate)
	r.Delete("/{domain}/templates/{name}", ms.deleteTemplate)
	r.Delete("/{domain}/templates/{name}", ms.deleteAllTemplates)

	ms.templates = append(ms.templates, mtypes.Template{
		Name:        "template1",
		Description: "template1 description",
		CreatedAt:   mtypes.RFC2822Time(time.Now()),
	})

	ms.templateVersions = make(map[string][]mtypes.TemplateVersion)
	ms.templateVersions["template1"] = []mtypes.TemplateVersion{
		{
			Tag:       "test",
			Template:  "template1 content",
			Engine:    "go",
			CreatedAt: mtypes.RFC2822Time(time.Now()),
			Comment:   "template1 comment",
			Active:    true,
		},
	}

	ms.templates = append(ms.templates, mtypes.Template{
		Name:        "template2",
		Description: "template2 description",
		CreatedAt:   mtypes.RFC2822Time(time.Now()),
	})

	ms.templateVersions["template2"] = []mtypes.TemplateVersion{
		{
			Tag:       "test",
			Template:  "template2 content",
			Engine:    "go",
			CreatedAt: mtypes.RFC2822Time(time.Now()),
			Comment:   "template2 comment",
			Active:    false,
		},
	}

	ms.templates = append(ms.templates, mtypes.Template{
		Name:        "template3",
		Description: "template3 description",
		CreatedAt:   mtypes.RFC2822Time(time.Now()),
	})

	ms.templateVersions["template3"] = []mtypes.TemplateVersion{
		{
			Tag:       "test",
			Template:  "template3 content",
			Engine:    "go",
			CreatedAt: mtypes.RFC2822Time(time.Now()),
			Comment:   "template3 comment",
			Active:    false,
		},
	}
}

func (ms *Server) listTemplates(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var idx []string
	for _, t := range ms.templates {
		idx = append(idx, t.Name)
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	page := r.FormValue("page")
	var pivot string
	if len(page) != 0 {
		pivot = r.FormValue("p")
		if pivot == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Invalid parameter: pivot \"}"))
			return
		}
	}
	start, end := pageOffsets(idx, page, pivot, limit)
	var nextAddress, prevAddress string
	var results []mtypes.Template

	if start != end {
		results = ms.templates[start:end]
		nextAddress = results[len(results)-1].Name
		prevAddress = results[0].Name
	} else {
		results = []mtypes.Template{}
		nextAddress = pivot
		prevAddress = pivot
	}

	toJSON(w, mtypes.ListTemplateResp{
		Paging: mtypes.Paging{
			First: getPageURL(r, url.Values{
				"page": []string{"first"},
			}),
			Last: getPageURL(r, url.Values{
				"page": []string{"last"},
			}),
			Next: getPageURL(r, url.Values{
				"page": []string{"next"},
				"p":    []string{nextAddress},
			}),
			Previous: getPageURL(r, url.Values{
				"page": []string{"prev"},
				"p":    []string{prevAddress},
			}),
		},
		Items: results,
	})
}

func (ms *Server) getTemplate(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "name")
	templateName = strings.ToLower(templateName)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	active := r.Form.Get("active")
	for _, template := range ms.templates {
		if template.Name == templateName {
			if active == "true" {
				version := ms.getActiveTemplateVersion(templateName)
				if version.Active { // active version exists
					template.Version = version
				}
			}
			toJSON(w, &mtypes.TemplateResp{
				Item: template,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"message\": \"template not found\"}"))
}

func (ms *Server) createTemplate(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	r.ParseForm()
	name := r.FormValue("name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"Missing mandatory parameter: name\"}"))
		return
	}

	name = strings.TrimSpace(name)

	if strings.Contains(name, " ") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"Invalid parameter: name is invalid\"}"))
		return
	}

	name = strings.ToLower(name)

	template := mtypes.Template{Name: name}
	template.CreatedAt = mtypes.RFC2822Time(time.Now())

	description := r.FormValue("description")
	if len(description) > 0 {
		template.Description = description
	}

	templateContent := r.FormValue("template")
	if len(templateContent) > 0 {
		templateVersion := mtypes.TemplateVersion{Template: templateContent}
		tag := r.FormValue("tag")
		if len(tag) > 0 {
			templateVersion.Tag = tag
		} else {
			templateVersion.Tag = "initial"
		}

		templateVersion.Comment = r.FormValue("comment")
		templateVersion.CreatedAt = mtypes.RFC2822Time(time.Now())
		templateVersion.Active = true

		engine := r.FormValue("engine")
		if len(engine) != 0 {
			if engine != "go" && engine != "handlebars" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("{\"message\": \"Invalid parameter: engine %s is not supported\"}", engine)))
				return
			}
			templateVersion.Engine = mtypes.TemplateEngine(engine)
		}

		template.Version = templateVersion
	}

	ms.templates = append(ms.templates, template)
	toJSON(w, map[string]any{
		"message":  "template has been stored",
		"template": template,
	})
}

func (ms *Server) updateTemplate(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"Missing mandatory parameter: name\"}"))
		return
	}
	name = strings.ToLower(name)

	r.ParseForm()
	description := r.FormValue("description")
	if len(description) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"No fields are provided to update\"}"))
		return
	}

	for i, template := range ms.templates {
		if template.Name == name {
			ms.templates[i].Description = description
			toJSON(w, map[string]any{
				"message": "template has been updated",
				"template": map[string]string{
					"name": name,
				},
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, map[string]any{
		"message": "template not found",
	})
}

func (ms *Server) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "name")

	for i, template := range ms.templates {
		if template.Name == templateName {
			ms.templates = append(ms.templates[:i], ms.templates[i+1:len(ms.templates)]...)

			toJSON(w, map[string]any{
				"message": "template has been deleted",
				"template": map[string]string{
					"name": templateName,
				},
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, map[string]string{"message": "template not found"})
}

func (ms *Server) deleteAllTemplates(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	ms.templates = []mtypes.Template{}
	ms.templateVersions = map[string][]mtypes.TemplateVersion{}

	toJSON(w, map[string]string{"message": "templates have been deleted"})
}

func (ms *Server) getActiveTemplateVersion(templateName string) mtypes.TemplateVersion {
	for _, templateVersion := range ms.templateVersions[templateName] {
		if templateVersion.Active {
			return templateVersion
		}
	}

	return mtypes.TemplateVersion{}
}
