package mailgun_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/pkg/errors"
)

func TestTemplateCRUD(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := mailgun.NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	findTemplate := func(id string) bool {
		it := mg.ListTemplates(nil)

		var page []mailgun.Template
		for it.Next(ctx, &page) {
			for _, template := range page {
				if template.Id == id {
					return true
				}
			}
		}
		ensure.Nil(t, it.Err())
		return false
	}

	const (
		Name        = "Mailgun-Go TestTemplateCRUD"
		Description = "Mailgun-Go Test Template Description"
		UpdatedDesc = "Mailgun-Go Test Updated Description"
	)

	tmpl := mailgun.Template{
		Name:        Name,
		Description: Description,
	}

	// Create a template
	ensure.Nil(t, mg.CreateTemplate(ctx, &tmpl))
	ensure.True(t, tmpl.Id != "")
	ensure.DeepEqual(t, tmpl.Description, Description)
	ensure.DeepEqual(t, tmpl.Name, Name)

	// Wait the template to show up
	ensure.Nil(t, waitForTemplate(mg, tmpl.Id))

	// Ensure the template is in the list
	ensure.True(t, findTemplate(tmpl.Id))

	// Update the description
	tmpl.Description = UpdatedDesc
	ensure.Nil(t, mg.UpdateTemplate(ctx, &tmpl))

	// Ensure update took
	updated, err := mg.GetTemplate(ctx, tmpl.Id)

	ensure.DeepEqual(t, updated.Id, tmpl.Id)
	ensure.DeepEqual(t, updated.Description, UpdatedDesc)

	// Delete the template
	ensure.Nil(t, mg.DeleteTemplate(ctx, tmpl.Id))
}

func waitForTemplate(mg mailgun.Mailgun, id string) error {
	ctx := context.Background()
	var attempts int
	for attempts <= 5 {
		_, err := mg.GetTemplate(ctx, id)
		if err != nil {
			if mailgun.GetStatusFromErr(err) == 404 {
				time.Sleep(time.Second * 2)
				attempts += 1
				continue
			}
			return err
		}
		return nil

	}
	return errors.Errorf("Waited to long for template '%s' to show up", id)
}
