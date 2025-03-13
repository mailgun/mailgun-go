package examples

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/events"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func AddBounce(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.AddBounce(ctx, domain, "bob@example.com", "550", "Undeliverable message error")
}

func CreateComplaint(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateComplaint(ctx, domain, "bob@example.com")
}

func AddDomain(domain, apiKey string) (mtypes.GetDomainResponse, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateDomain(ctx, domain, &mailgun.CreateDomainOptions{
		Password:   "super_secret",
		SpamAction: mtypes.SpamActionTag,
		Wildcard:   false,
	})
}

func AddDomainIPS(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.AddDomainIP(ctx, domain, "127.0.0.1")
}

func AddListMember(apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	memberJoe := mtypes.Member{
		Address:    "joe@example.com",
		Name:       "Joe Example",
		Subscribed: mtypes.Subscribed,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "mailingList@example.com", memberJoe)
}

func AddListMembers(apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMemberList(ctx, nil, "mailgunList@example.com", []any{
		mtypes.Member{
			Address:    "alice@example.com",
			Name:       "Alice's debugging account",
			Subscribed: mtypes.Unsubscribed,
		},
		mtypes.Member{
			Address:    "Bob Cool <bob@example.com>",
			Name:       "Bob's Cool Account",
			Subscribed: mtypes.Subscribed,
		},
		mtypes.Member{
			Address: "joe.hamradio@example.com",
			// Charlette is a ham radio packet BBS user.
			// We attach her packet BBS e-mail address as an arbitrary var here.
			Vars: map[string]any{
				"packet-email": "KW9ABC @ BOGUS-4.#NCA.CA.USA.NOAM",
			},
		},
	})
}

func CreateUnsubscribe(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateUnsubscribe(ctx, domain, "bob@example.com", "*")
}

func CreateUnsubscribes(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	unsubscribes := []mtypes.Unsubscribe{
		{Address: "alice@example.com"},
		{Address: "bob@example.com", Tags: []string{"tag1"}},
	}

	return mg.CreateUnsubscribes(ctx, domain, unsubscribes)
}

func CreateUnsubscribeWithTag(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateUnsubscribe(ctx, domain, "bob@example.com", "tag1")
}

func CreateWebhook(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateWebhook(ctx, domain, "clicked", []string{"https://your_domain.com/v1/clicked"})
}

func ChangePassword(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ChangeCredentialPassword(ctx, domain, "alice", "super_secret")
}

func CreateCredential(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateCredential(ctx, domain, "alice@example.com", "secret")
}

func CreateDomain(domain, apiKey string) (mtypes.GetDomainResponse, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateDomain(ctx, domain, &mailgun.CreateDomainOptions{
		Password:   "super_secret",
		SpamAction: mtypes.SpamActionTag,
		Wildcard:   false,
	})
}

func CreateExport(apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateExport(ctx, "/v3/domains")
}

func CreateMailingList(apiKey string) (mtypes.MailingList, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMailingList(ctx, mtypes.MailingList{
		Address:     "list@example.com",
		Name:        "dev",
		Description: "Mailgun developers list.",
		AccessLevel: mtypes.AccessLevelMembers,
	})
}

func CreateRoute(domain, apiKey string) (mtypes.Route, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateRoute(ctx, mtypes.Route{
		Priority:    1,
		Description: "Sample Route",
		Expression:  "match_recipient(\".*@YOUR_DOMAIN_NAME\")",
		Actions: []string{
			"forward(\"http://example.com/messages/\")",
			"stop()",
		},
	})
}

func DeleteCredential(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteCredential(ctx, domain, "alice")
}

func DeleteDomain(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteDomain(ctx, domain)
}

func DeleteTag(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteTag(ctx, domain, "newsletter")
}

func DeleteWebhook(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteWebhook(ctx, domain, "clicked")
}

func PrintEventLog(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	// Create an iterator
	it := mg.ListEvents(domain, &mailgun.ListEventOptions{
		Begin: time.Now().Add(-50 * time.Minute),
		Limit: 100,
		Filter: map[string]string{
			"recipient": "joe@example.com",
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Iterate through all the pages of events
	var page []events.Event
	for it.Next(ctx, &page) {
		for _, event := range page {
			fmt.Printf("%+v\n", event)
		}
	}

	// Did iteration end because of an error?
	if it.Err() != nil {
		return it.Err()
	}

	return nil
}

func PrintFailedEvents(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	// Create an iterator
	it := mg.ListEvents(domain, &mailgun.ListEventOptions{
		Filter: map[string]string{
			"event": "rejected OR failed",
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Iterate through all the pages of events
	var page []events.Event
	for it.Next(ctx, &page) {
		for _, event := range page {
			switch e := event.(type) {
			case *events.Failed:
				fmt.Printf("Failed Reason: %s", e.Reason)
			case *events.Rejected:
				fmt.Printf("Rejected Reason: %s", e.Reject.Reason)
			}
		}
	}

	// Did iteration end because of an error?
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}

func PrintEvents(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	// Create an iterator
	it := mg.ListEvents(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Iterate through all the pages of events
	var page []events.Event
	for it.Next(ctx, &page) {
		for _, event := range page {
			switch e := event.(type) {
			case *events.Accepted:
				fmt.Printf("Accepted ID: %s", e.Message.Headers.MessageID)
			case *events.Rejected:
				fmt.Printf("Rejected Reason: %s", e.Reject.Reason)
				// Add other event types here
			}
			fmt.Printf("%+v\n", event.GetTimestamp())
		}
	}

	// Did iteration end because of an error?
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}

func GetBounce(domain, apiKey string) (mtypes.Bounce, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetBounce(ctx, domain, "foo@bar.com")
}

func ListBounces(domain, apiKey string) ([]mtypes.Bounce, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListBounces(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Bounce
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetComplaints(domain, apiKey string) (mtypes.Complaint, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetComplaint(ctx, domain, "baz@example.com")
}

func ListComplaints(domain, apiKey string) ([]mtypes.Complaint, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListComplaints(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Complaint
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetDomainConnection(domain, apiKey string) (mtypes.DomainConnection, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetDomainConnection(ctx, domain)
}

func ListCredentials(domain, apiKey string) ([]mtypes.Credential, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListCredentials(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Credential
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetDomain(domain, apiKey string) (mtypes.GetDomainResponse, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetDomain(ctx, domain, nil)
}

func ListDomainIPS(domain, apiKey string) ([]mtypes.IPAddress, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ListDomainIPS(ctx, domain)
}

func GetDomainTracking(domain, apiKey string) (mtypes.DomainTracking, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetDomainTracking(ctx, domain)
}

func ListDomains(domain, apiKey string) ([]mtypes.Domain, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListDomains(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Domain
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetExport(domain, apiKey string) (mtypes.Export, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetExport(ctx, "EXPORT_ID")
}

func GetIP(domain, apiKey string) (mtypes.IPAddress, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetIP(ctx, "127.0.0.1")
}

func ListIPS(domain, apiKey string) ([]mtypes.IPAddress, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Pass 'true' to only return dedicated ips
	return mg.ListIPS(ctx, true)
}

func GetTagLimits(domain, apiKey string) (mtypes.TagLimits, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetTagLimits(ctx, domain)
}

func ListExports(apiKey string) ([]mtypes.Export, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Optionally pass a url to filter by
	return mg.ListExports(ctx, "")
}

func GetMembers(apiKey string) ([]mtypes.Member, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListMembers("list@example.com", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Member
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListMailingLists(apiKey string) ([]mtypes.MailingList, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListMailingLists(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.MailingList
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetRoute(domain, apiKey string) (mtypes.Route, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetRoute(ctx, "route_id")
}

func ListRoutes(domain, apiKey string) ([]mtypes.Route, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListRoutes(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Route
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListTags(domain, apiKey string) ([]mtypes.Tag, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListTags(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Tag
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListUnsubscribes(domain, apiKey string) ([]mtypes.Unsubscribe, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListUnsubscribes(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Unsubscribe
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ValidateEmail(apiKey string) (mtypes.ValidateEmailResponse, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ValidateEmail(ctx, "foo@mailgun.net", false)
}

func GetWebhook(domain, apiKey string) ([]string, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetWebhook(ctx, domain, "clicked")
}

func ListWebhooks(domain, apiKey string) (map[string][]string, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ListWebhooks(ctx, domain)
}

func DeleteDomainIP(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteDomainIP(ctx, domain, "127.0.0.1")
}

func DeleteListMember(apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteMember(ctx, "joe@example.com", "list@example.com")
}

func DeleteMailingList(apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteMailingList(ctx, "list@example.com")
}

func ResendMessage(domain, apiKey string) (mtypes.SendMessageResponse, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ReSend(ctx, "STORAGE_URL", "bar@example.com")
}

func SendComplexMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"foo@example.com",
	)
	m.AddCC("baz@example.com")
	m.AddBCC("bar@example.com")
	m.SetHTML("<html>HTML version of the body</html>")
	m.AddAttachment("files/test.jpg")
	m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendWithConnectionOptions(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"foo@example.com",
	)

	m.SetRequireTLS(true)
	m.SetSkipVerification(true)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendInlineImage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"foo@example.com",
	)
	m.AddCC("baz@example.com")
	m.AddBCC("bar@example.com")
	m.SetHTML(`<html>Inline image here: <img alt="image" src="cid:test.jpg"/></html>`)
	m.AddInline("files/test.jpg")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendMessageNoTracking(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"foo@example.com",
	)
	m.SetTracking(false)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendMimeMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	mimeMsgReader, err := os.Open("files/message.mime")
	if err != nil {
		return "", err
	}

	m := mailgun.NewMIMEMessage(domain, mimeMsgReader, "bar@example.com")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendScheduledMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"bar@example.com",
	)
	m.SetDeliveryTime(time.Now().Add(5 * time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendSimpleMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <mailgun@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"YOU@YOUR_DOMAIN_NAME",
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendTaggedMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"bar@example.com",
	)

	err := m.AddTag("FooTag", "BarTag", "BlortTag")
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendTemplateMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(apiKey)
	m := mailgun.NewMessage(
		domain,
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hey %recipient.first%",
		"If you wish to unsubscribe, click http://mailgun/unsubscribe/%recipient.id%",
	) // IMPORTANT: No To:-field recipients!

	// Set template to be applied to this message.
	m.SetTemplate("my-template")

	m.AddRecipientAndVariables("bob@example.com", map[string]any{
		"first": "bob",
		"id":    1,
	})

	m.AddRecipientAndVariables("alice@example.com", map[string]any{
		"first": "alice",
		"id":    2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func UpdateDomainConnection(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateDomainConnection(ctx, domain, mtypes.DomainConnection{
		RequireTLS:       true,
		SkipVerification: true,
	})
}

func UpdateMember(apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := mg.UpdateMember(ctx, "bar@example.com", "list@example.com", mtypes.Member{
		Name:       "Foo Bar",
		Subscribed: mtypes.Unsubscribed,
	})
	return err
}

func UpdateWebhook(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateWebhook(ctx, domain, "clicked", []string{"https://your_domain.com/clicked"})
}

func VerifyWebhookSignature(apiKey, webhookSigningKey, timestamp, token, signature string) (bool, error) {
	mg := mailgun.NewMailgun(apiKey)
	mg.SetWebhookSigningKey(webhookSigningKey)

	return mg.VerifyWebhookSignature(mtypes.Signature{
		TimeStamp: timestamp,
		Token:     token,
		Signature: signature,
	})
}

func SendMessageWithTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create a new template
	err = mg.CreateTemplate(ctx, domain, &mtypes.Template{
		Name: "my-template",
		Version: mtypes.TemplateVersion{
			Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
			Engine:   mtypes.TemplateEngineGo,
			Tag:      "v1",
		},
	})
	if err != nil {
		return err
	}

	// Give time for template to show up in the system.
	time.Sleep(time.Second * 1)

	// Create a new message with template
	m := mailgun.NewMessage(domain, "Excited User <excited@example.com>", "Template example", "")
	m.SetTemplate("my-template")

	// Add recipients
	m.AddRecipient("bob@example.com")
	m.AddRecipient("alice@example.com")

	// Add the variables to be used by the template
	m.AddVariable("title", "Hello Templates")
	m.AddVariable("body", "Body of the message")

	_, id, err := mg.Send(ctx, m)
	fmt.Printf("Queued: %s", id)
	return err
}

func CreateTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateTemplate(ctx, domain, &mtypes.Template{
		Name: "my-template",
		Version: mtypes.TemplateVersion{
			Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
			Engine:   mtypes.TemplateEngineGo,
			Tag:      "v1",
		},
	})
}

func DeleteTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteTemplate(ctx, domain, "my-template")
}

func UpdateTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateTemplate(ctx, domain, &mtypes.Template{
		Name:        "my-template",
		Description: "Add a description to the template",
	})
}

func GetTemplate(domain, apiKey string) (mtypes.Template, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetTemplate(ctx, domain, "my-template")
}

func ListActiveTemplates(domain, apiKey string) ([]mtypes.Template, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListTemplates(domain, &mailgun.ListTemplateOptions{Active: true})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Template
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListTemplates(domain, apiKey string) ([]mtypes.Template, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListTemplates(domain, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.Template
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func AddTemplateVersion(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.AddTemplateVersion(ctx, domain, "my-template", &mtypes.TemplateVersion{
		Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
		Engine:   mtypes.TemplateEngineGo,
		Tag:      "v2",
		Active:   true,
	})
}

func DeleteTemplateVersion(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Delete the template version tagged as 'v2'
	return mg.DeleteTemplateVersion(ctx, domain, "my-template", "v2")
}

func GetTemplateVersion(domain, apiKey string) (mtypes.TemplateVersion, error) {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Get the template version tagged as 'v2'
	return mg.GetTemplateVersion(ctx, domain, "my-template", "v2")
}

func UpdateTemplateVersion(domain, apiKey string) error {
	mg := mailgun.NewMailgun(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateTemplateVersion(ctx, domain, "my-template", &mtypes.TemplateVersion{
		Comment: "Add a comment to the template and make it 'active'",
		Tag:     "v2",
		Active:  true,
	})
}

func ListTemplateVersions(domain, apiKey string) ([]mtypes.TemplateVersion, error) {
	mg := mailgun.NewMailgun(apiKey)
	it := mg.ListTemplateVersions(domain, "my-template", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mtypes.TemplateVersion
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}
