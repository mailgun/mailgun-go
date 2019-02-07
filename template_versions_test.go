package mailgun_test

import (
	"context"
	"fmt"
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

	findVersion := func(templateName, tag string) bool {
		fmt.Printf("List\n")
		it := mg.ListTemplateVersions(templateName, nil)

		var page []mailgun.TemplateVersion
		for it.Next(ctx, &page) {
			for _, v := range page {
				if v.Tag == tag {
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
		Tag            = "v1"
	)

	tmpl := mailgun.Template{
		Name: "Mailgun-go-TestTemplateVersionsCRUD",
	}

	// Create a template
	ensure.Nil(t, mg.CreateTemplate(ctx, &tmpl))

	version := mailgun.TemplateVersion{
		Tag:      Tag,
		Comment:  Comment,
		Template: Template,
		Active:   true,
		Engine:   mailgun.TemplateEngineGo,
	}

	// Add a version version
	ensure.Nil(t, mg.AddTemplateVersion(ctx, tmpl.Name, &version))
	ensure.DeepEqual(t, version.Tag, Tag)
	ensure.DeepEqual(t, version.Comment, Comment)
	ensure.DeepEqual(t, version.Engine, mailgun.TemplateEngineGo)

	// Ensure the version is in the list
	ensure.True(t, findVersion(tmpl.Name, version.Tag))

	// Update the Comment
	version.Comment = UpdatedComment
	ensure.Nil(t, mg.UpdateTemplateVersion(ctx, tmpl.Name, &version))

	// Ensure update took
	updated, err := mg.GetTemplateVersion(ctx, tmpl.Name, version.Tag)

	ensure.DeepEqual(t, updated.Comment, UpdatedComment)

	// Add a new active Version
	version2 := mailgun.TemplateVersion{
		Tag:      "v2",
		Comment:  Comment,
		Template: Template,
		Active:   true,
		Engine:   mailgun.TemplateEngineGo,
	}
	ensure.Nil(t, mg.AddTemplateVersion(ctx, tmpl.Name, &version2))

	// Ensure the version is in the list
	ensure.True(t, findVersion(tmpl.Name, version2.Tag))

	// Delete the first version
	ensure.Nil(t, mg.DeleteTemplateVersion(ctx, tmpl.Name, version.Tag))

	// Ensure version was deleted
	ensure.False(t, findVersion(tmpl.Name, version.Tag))

	// Delete the template
	ensure.Nil(t, mg.DeleteTemplate(ctx, tmpl.Name))
}
