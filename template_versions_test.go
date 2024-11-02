package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestTemplateVersionsCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	findVersion := func(templateName, tag string) bool {
		it := mg.ListTemplateVersions(templateName, nil)

		var page []mailgun.TemplateVersion
		for it.Next(ctx, &page) {
			for _, v := range page {
				if v.Tag == tag {
					return true
				}
			}
		}
		require.NoError(t, it.Err())
		return false
	}

	const (
		Comment        = "Mailgun-Go TestTemplateVersionsCRUD"
		UpdatedComment = "Mailgun-Go Test Version Updated"
		Template       = "{{.Name}}"
		Tag            = "v1"
	)

	tmpl := mailgun.Template{
		Name: randomString(10, "Mailgun-go-TestTemplateVersionsCRUD-"),
	}

	// Create a template
	require.NoError(t, mg.CreateTemplate(ctx, &tmpl))

	version := mailgun.TemplateVersion{
		Tag:      Tag,
		Comment:  Comment,
		Template: Template,
		Active:   true,
		Engine:   mailgun.TemplateEngineGo,
	}

	// Add a version version
	require.NoError(t, mg.AddTemplateVersion(ctx, tmpl.Name, &version))
	ensure.DeepEqual(t, version.Tag, Tag)
	ensure.DeepEqual(t, version.Comment, Comment)
	ensure.DeepEqual(t, version.Engine, mailgun.TemplateEngineGo)

	// Ensure the version is in the list
	ensure.True(t, findVersion(tmpl.Name, version.Tag))

	// Update the Comment
	version.Comment = UpdatedComment
	version.Template = Template + "updated"
	require.NoError(t, mg.UpdateTemplateVersion(ctx, tmpl.Name, &version))

	// Ensure update took
	updated, err := mg.GetTemplateVersion(ctx, tmpl.Name, version.Tag)

	require.NoError(t, err)
	ensure.DeepEqual(t, updated.Comment, UpdatedComment)
	ensure.DeepEqual(t, updated.Template, Template+"updated")

	// Add a new active Version
	version2 := mailgun.TemplateVersion{
		Tag:      "v2",
		Comment:  Comment,
		Template: Template,
		Active:   true,
		Engine:   mailgun.TemplateEngineGo,
	}
	require.NoError(t, mg.AddTemplateVersion(ctx, tmpl.Name, &version2))

	// Ensure the version is in the list
	ensure.True(t, findVersion(tmpl.Name, version2.Tag))

	// Delete the first version
	require.NoError(t, mg.DeleteTemplateVersion(ctx, tmpl.Name, version.Tag))

	// Ensure version was deleted
	ensure.False(t, findVersion(tmpl.Name, version.Tag))

	// Delete the template
	require.NoError(t, mg.DeleteTemplate(ctx, tmpl.Name))
}
