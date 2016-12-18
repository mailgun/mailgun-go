package mailgun

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestWebhookCRUD(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	var countHooks = func() int {
		hooks, err := mg.GetWebhooks()
		ensure.Nil(t, err)
		return len(hooks)
	}

	hookCount := countHooks()

	ensure.Nil(t, mg.CreateWebhook("deliver", "http://www.example.com"))
	defer func() {
		ensure.Nil(t, mg.DeleteWebhook("deliver"))
		newCount := countHooks()
		ensure.DeepEqual(t, newCount, hookCount)
	}()

	newCount := countHooks()
	ensure.False(t, newCount <= hookCount)

	theURL, err := mg.GetWebhookByType("deliver")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theURL, "http://www.example.com")

	ensure.Nil(t, mg.UpdateWebhook("deliver", "http://api.example.com"))

	hooks, err := mg.GetWebhooks()
	ensure.Nil(t, err)

	ensure.DeepEqual(t, hooks["deliver"], "http://api.example.com")
}
