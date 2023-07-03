package mailgun

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addTemplateVersionRoutes(r chi.Router) {
	r.Get("/{domain}/templates/{template}/versions", ms.listTemplateVersions)
	r.Get("/{domain}/templates/{template}/versions/{tag}", ms.getTemplateVersion)
	r.Post("/{domain}/templates/{template}/versions", ms.createTemplateVersion)
	r.Put("/{domain}/templates/{template}/versions/{tag}", ms.updateTemplateVersion)
	r.Delete("/{domain}/templates/{template}/versions/{tag}", ms.deleteTemplateVersion)
}

func (ms *mockServer) listTemplateVersions(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "template")
	template, found := ms.fetchTemplate(templateName)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template not found\"}"))
		return
	}

	templateVersions := ms.templateVersions[templateName]

	var idx []string
	for _, t := range templateVersions {
		idx = append(idx, t.Tag)
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 10
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
	var results []TemplateVersion

	if start != end {
		results = ms.templateVersions[templateName][start:end]
		nextAddress = results[len(results)-1].Tag
		prevAddress = results[0].Tag
	} else {
		results = []TemplateVersion{}
		nextAddress = pivot
		prevAddress = pivot
	}

	toJSON(w, templateVersionListResp{
		Paging: Paging{
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
		Template: struct {
			Template
			Versions []TemplateVersion `json:"versions,omitempty"`
		}{
			Template: template,
			Versions: results,
		},
	})
}

func (ms *mockServer) getTemplateVersion(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "template")
	templateName = strings.ToLower(templateName)
	templateVersionName := chi.URLParam(r, "tag")
	templateVersionName = strings.ToLower(templateVersionName)

	template, templateFound := ms.fetchTemplate(templateName)
	if !templateFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template not found\"}"))
		return
	}

	templateVersion, templateVersionFound := ms.fetchTemplateVersion(templateName, templateVersionName)
	if !templateVersionFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template version not found\"}"))
		return
	}

	template.Version = templateVersion

	toJSON(w, &templateResp{
		Item: template,
	})
}

func (ms *mockServer) createTemplateVersion(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "template")
	templateName = strings.ToLower(templateName)

	r.ParseForm()
	templateContent := r.FormValue("template")
	if len(templateContent) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"Missing mandatory parameter: template\"}"))
		return
	}
	tagName := r.FormValue("tag")
	if len(templateContent) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"Missing mandatory parameter: tag\"}"))
		return
	}

	template, templateFound := ms.fetchTemplate(templateName)
	if !templateFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template not found\"}"))
		return
	}

	_, templateVersionFound := ms.fetchTemplateVersion(templateName, tagName)
	if templateVersionFound {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"version %s already exists\"}", tagName)))
		return
	}

	comment := r.FormValue("comment")
	active := r.FormValue("active")

	engine := r.FormValue("engine")
	if len(engine) != 0 {
		if strings.ToLower(engine) != "go" && strings.ToLower(engine) != "handlebars" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"message\": \"Invalid parameter: engine %s is not supported\"}", engine)))
			return
		}
	} else {
		engine = "handlebars"
	}

	newTemplateVersion := TemplateVersion{
		Template:  templateContent,
		Comment:   comment,
		Tag:       tagName,
		Engine:    TemplateEngine(engine),
		CreatedAt: RFC2822Time(time.Now()),
	}

	if active == "yes" {
		newTemplateVersion.Active = true
		for i, _ := range ms.templateVersions[templateName] {
			ms.templateVersions[templateName][i].Active = false
		}
	}

	ms.templateVersions[templateName] = append(ms.templateVersions[templateName], newTemplateVersion)
	template.Version = newTemplateVersion
	toJSON(w, map[string]interface{}{
		"message":  "new version of the template has been stored",
		"template": template,
	})
}

func (ms *mockServer) updateTemplateVersion(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "template")
	templateName = strings.ToLower(templateName)
	templateVersionName := chi.URLParam(r, "tag")
	templateVersionName = strings.ToLower(templateVersionName)

	_, templateFound := ms.fetchTemplate(templateName)
	if !templateFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template not found\"}"))
		return
	}

	var templateVersionFound bool
	var templateVersion TemplateVersion
	var templateVersionIndex int
	for i, tmplVersion := range ms.templateVersions[templateName] {
		if tmplVersion.Tag == templateVersionName {
			templateVersion = tmplVersion
			templateVersionFound = true
			templateVersionIndex = i
			break
		}
	}

	if !templateVersionFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template version not found\"}"))
		return
	}

	r.ParseForm()
	templateContent := r.FormValue("template")
	comment := r.FormValue("comment")
	active := r.FormValue("active")

	updated := false
	if len(templateContent) != 0 {
		templateVersion.Template = templateContent
		updated = true
	}
	if len(comment) != 0 {
		templateVersion.Comment = comment
		updated = true
	}
	if len(active) != 0 {
		if active == "yes" {
			templateVersion.Active = true
			for i := range ms.templateVersions[templateName] { //every other template version become not active
				if i == templateVersionIndex {
					continue
				}
				ms.templateVersions[templateName][i].Active = false
			}
		}
		updated = true
	}

	if !updated {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"No fields are provided to update\"}"))
		return
	}

	ms.templateVersions[templateName][templateVersionIndex] = templateVersion
	toJSON(w, map[string]interface{}{
		"message": "version has been updated",
		"template": map[string]interface{}{
			"name": templateName,
			"version": map[string]string{
				"tag": templateVersionName,
			},
		},
	})
}

func (ms *mockServer) deleteTemplateVersion(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	templateName := chi.URLParam(r, "template")
	templateName = strings.ToLower(templateName)
	templateVersionName := chi.URLParam(r, "tag")
	templateVersionName = strings.ToLower(templateVersionName)

	_, templateFound := ms.fetchTemplate(templateName)
	if !templateFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"template not found\"}"))
		return
	}
	for i, templateVersion := range ms.templateVersions[templateName] {
		if templateVersion.Tag == templateVersionName {
			ms.templateVersions[templateName] = append(ms.templateVersions[templateName][:i], ms.templateVersions[templateName][i+1:len(ms.templateVersions[templateName])]...)
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	toJSON(w, map[string]interface{}{
		"message": "version has been deleted",
		"template": map[string]interface{}{
			"name": templateName,
			"version": map[string]string{
				"tag": templateVersionName,
			},
		},
	})
}

func (ms *mockServer) fetchTemplate(name string) (template Template, found bool) {
	for _, existingTemplate := range ms.templates {
		if existingTemplate.Name == name {
			template = existingTemplate
			return template, true
		}
	}

	return Template{}, false
}

func (ms *mockServer) fetchTemplateVersion(templateName string, templateVersionTag string) (TemplateVersion, bool) {
	for _, existingTemplate := range ms.templateVersions[templateName] {
		if existingTemplate.Tag == templateVersionTag {
			return existingTemplate, true
		}
	}

	return TemplateVersion{}, false
}
