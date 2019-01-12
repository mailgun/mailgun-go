package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3"
)

func TestTemplateVersionsCRUD(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := mailgun.NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	findVersion := func(templateId, versionId string) bool {
		it := mg.ListTemplateVersions(templateId, nil)

		var page []mailgun.TemplateVersion
		for it.Next(ctx, &page) {
			for _, v := range page {
				if v.Id == versionId {
					return true
				}
			}
		}
		ensure.Nil(t, it.Err())
		return false
	}

	const (
		Comment        = "Mailgun-Go TestTemplateVersionsCRUD"
		UpdatedComment = "Mailgun-Go Test Version Updated"
		Template       = "{{.Name}}"
	)

	tmpl := mailgun.Template{
		Name: "Mailgun-go TestTemplateVersionsCRUD",
	}

	// Create a template
	ensure.Nil(t, mg.CreateTemplate(ctx, &tmpl))

	version := mailgun.TemplateVersion{
		Comment:  Comment,
		Template: Template,
		Active:   true,
		Engine:   mailgun.TemplateEngineGo,
	}

	// Add a version version
	ensure.Nil(t, mg.AddTemplateVersion(ctx, tmpl.Id, &version))
	ensure.True(t, version.Id != "")
	ensure.DeepEqual(t, version.Comment, Comment)
	ensure.DeepEqual(t, version.Engine, mailgun.TemplateEngineGo)

	// Ensure the version is in the list
	ensure.True(t, findVersion(tmpl.Id, version.Id))

	// Update the Comment
	version.Comment = UpdatedComment
	ensure.Nil(t, mg.UpdateTemplateVersion(ctx, tmpl.Id, &version))

	// Ensure update took
	updated, err := mg.GetTemplateVersion(ctx, tmpl.Id, version.Id)

	ensure.DeepEqual(t, version.Id, updated.Id)
	ensure.DeepEqual(t, updated.Comment, UpdatedComment)

	// Add a new active Version
	version2 := mailgun.TemplateVersion{
		Comment:  Comment,
		Template: Template,
		Active:   true,
		Engine:   mailgun.TemplateEngineGo,
	}
	ensure.Nil(t, mg.AddTemplateVersion(ctx, tmpl.Id, &version2))

	// Ensure the version is in the list
	ensure.True(t, findVersion(tmpl.Id, version2.Id))

	// Delete the first version
	ensure.Nil(t, mg.DeleteTemplateVersion(ctx, tmpl.Id, version.Id))

	// Ensure version was deleted
	ensure.False(t, findVersion(tmpl.Id, version.Id))

	// Delete the template
	ensure.Nil(t, mg.DeleteTemplate(ctx, tmpl.Id))
}
