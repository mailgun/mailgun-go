package examples

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/mailgun/mailgun-go/v3/events"
	"os"
	"time"
)

func AddBounce(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.AddBounce(ctx, "bob@example.com", "550", "Undeliverable message error")
}

func CreateComplaint(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateComplaint(ctx, "bob@example.com")
}

func AddDomain(domain, apiKey string) (mailgun.DomainResponse, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateDomain(ctx, "example.com", &mailgun.CreateDomainOptions{
		Password:   "super_secret",
		SpamAction: mailgun.SpamActionTag,
		Wildcard:   false,
	})
}

func AddDomainIPS(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.AddDomainIP(ctx, "127.0.0.1")
}

func AddListMember(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	memberJoe := mailgun.Member{
		Address:    "joe@example.com",
		Name:       "Joe Example",
		Subscribed: mailgun.Subscribed,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "mailingList@example.com", memberJoe)
}

func AddListMembers(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMemberList(ctx, nil, "mailgunList@example.com", []interface{}{
		mailgun.Member{
			Address:    "alice@example.com",
			Name:       "Alice's debugging account",
			Subscribed: mailgun.Unsubscribed,
		},
		mailgun.Member{
			Address:    "Bob Cool <bob@example.com>",
			Name:       "Bob's Cool Account",
			Subscribed: mailgun.Subscribed,
		},
		mailgun.Member{
			Address: "joe.hamradio@example.com",
			// Charlette is a ham radio packet BBS user.
			// We attach her packet BBS e-mail address as an arbitrary var here.
			Vars: map[string]interface{}{
				"packet-email": "KW9ABC @ BOGUS-4.#NCA.CA.USA.NOAM",
			},
		},
	})
}

func CreateUnsubscribe(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateUnsubscribe(ctx, "bob@example.com", "*")
}

func CreateUnsubscribeWithTag(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateUnsubscribe(ctx, "bob@example.com", "tag1")
}

func CreateWebhook(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateWebhook(ctx, "clicked", []string{"https://your_domain.com/v1/clicked"})
}

func ChangePassword(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ChangeCredentialPassword(ctx, "alice", "super_secret")
}

func CreateCredential(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateCredential(ctx, "alice@example.com", "secret")
}

func CreateDomain(domain, apiKey string) (mailgun.DomainResponse, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateDomain(ctx, "example.com", &mailgun.CreateDomainOptions{
		Password:   "super_secret",
		SpamAction: mailgun.SpamActionTag,
		Wildcard:   false,
	})
}

func CreateExport(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateExport(ctx, "/v3/domains")
}

func CreateMailingList(domain, apiKey string) (mailgun.MailingList, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMailingList(ctx, mailgun.MailingList{
		Address:     "list@example.com",
		Name:        "dev",
		Description: "Mailgun developers list.",
		AccessLevel: mailgun.AccessLevelMembers,
	})
}

func CreateRoute(domain, apiKey string) (mailgun.Route, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateRoute(ctx, mailgun.Route{
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
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteCredential(ctx, "alice")
}

func DeleteDomain(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteDomain(ctx, "example.com")
}

func DeleteTag(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteTag(ctx, "newsletter")
}

func DeleteWebhook(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteWebhook(ctx, "clicked")
}

func PrintEventLog(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	// Create an iterator
	it := mg.ListEvents(&mailgun.ListEventOptions{
		Begin: time.Now().Add(-50 * time.Minute),
		Limit: 100,
		Filter: map[string]string{
			"recipient": "joe@example.com",
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Iterate through all the pages of events
	var page []mailgun.Event
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
	mg := mailgun.NewMailgun(domain, apiKey)

	// Create an iterator
	it := mg.ListEvents(&mailgun.ListEventOptions{
		Filter: map[string]string{
			"event": "rejected OR failed",
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Iterate through all the pages of events
	var page []mailgun.Event
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
	mg := mailgun.NewMailgun(domain, apiKey)

	// Create an iterator
	it := mg.ListEvents(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Iterate through all the pages of events
	var page []mailgun.Event
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

func GetBounce(domain, apiKey string) (mailgun.Bounce, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetBounce(ctx, "foo@bar.com")
}

func ListBounces(domain, apiKey string) ([]mailgun.Bounce, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListBounces(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Bounce
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetComplaints(domain, apiKey string) (mailgun.Complaint, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetComplaint(ctx, "baz@example.com")
}

func ListComplaints(domain, apiKey string) ([]mailgun.Complaint, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListComplaints(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Complaint
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetDomainConnection(domain, apiKey string) (mailgun.DomainConnection, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetDomainConnection(ctx, domain)
}

func ListCredentials(domain, apiKey string) ([]mailgun.Credential, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListCredentials(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Credential
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetDomain(domain, apiKey string) (mailgun.DomainResponse, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetDomain(ctx, domain)
}

func ListDomainIPS(domain, apiKey string) ([]mailgun.IPAddress, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ListDomainIPS(ctx)
}

func GetDomainTracking(domain, apiKey string) (mailgun.DomainTracking, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetDomainTracking(ctx, domain)
}

func ListDomains(domain, apiKey string) ([]mailgun.Domain, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListDomains(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Domain
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetExport(domain, apiKey string) (mailgun.Export, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetExport(ctx, "EXPORT_ID")
}

func GetIP(domain, apiKey string) (mailgun.IPAddress, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetIP(ctx, "127.0.0.1")
}

func ListIPS(domain, apiKey string) ([]mailgun.IPAddress, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Pass 'true' to only return dedicated ips
	return mg.ListIPS(ctx, true)
}

func GetTagLimits(domain, apiKey string) (mailgun.TagLimits, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetTagLimits(ctx, domain)
}

func ListExports(domain, apiKey string) ([]mailgun.Export, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Optionally pass a url to filter by
	return mg.ListExports(ctx, "")
}

func GetMembers(domain, apiKey string) ([]mailgun.Member, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListMembers("list@example.com", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Member
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListMailingLists(domain, apiKey string) ([]mailgun.MailingList, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListMailingLists(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.MailingList
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ParseAddress(apiKey string) ([]string, []string, error) {
	mv := mailgun.NewEmailValidator(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mv.ParseAddresses(ctx,
		"Alice <alice@example.com>",
		"bob@example.com",
		// ...
	)
}

func GetRoute(domain, apiKey string) (mailgun.Route, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetRoute(ctx, "route_id")
}

func ListRoutes(domain, apiKey string) ([]mailgun.Route, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListRoutes(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Route
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func GetStats(domain, apiKey string) ([]mailgun.Stats, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetStats(ctx, []string{"accepted", "delivered", "failed"}, &mailgun.GetStatOptions{
		Duration: "1m",
	})
}

func ListTags(domain, apiKey string) ([]mailgun.Tag, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListTags(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Tag
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListUnsubscribes(domain, apiKey string) ([]mailgun.Unsubscribe, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListUnsubscribes(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Unsubscribe
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ValidateEmail(apiKey string) (mailgun.EmailVerification, error) {
	mv := mailgun.NewEmailValidator(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mv.ValidateEmail(ctx, "foo@mailgun.net", false)
}

func GetWebhook(domain, apiKey string) ([]string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetWebhook(ctx, "clicked")
}

func ListWebhooks(domain, apiKey string) (map[string][]string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ListWebhooks(ctx)
}

func DeleteDomainIP(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteDomainIP(ctx, "127.0.0.1")
}

func DeleteListMember(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteMember(ctx, "joe@example.com", "list@example.com")
}

func DeleteMailingList(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteMailingList(ctx, "list@example.com")
}

func ResendMessage(domain, apiKey string) (string, string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.ReSend(ctx, "STORAGE_URL", "bar@example.com")
}

func SendComplexMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"foo@example.com",
	)
	m.AddCC("baz@example.com")
	m.AddBCC("bar@example.com")
	m.SetHtml("<html>HTML version of the body</html>")
	m.AddAttachment("files/test.jpg")
	m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendWithConnectionOptions(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
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
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"foo@example.com",
	)
	m.AddCC("baz@example.com")
	m.AddBCC("bar@example.com")
	m.SetHtml("<html>HTML version of the body</html>")
	m.AddInline("files/test.jpg")
	m.AddInline("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendMessageNoTracking(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
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
	mg := mailgun.NewMailgun(domain, apiKey)
	mimeMsgReader, err := os.Open("files/message.mime")
	if err != nil {
		return "", err
	}

	m := mg.NewMIMEMessage(mimeMsgReader, "bar@example.com")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func SendScheduledMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
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
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
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
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
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
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Excited User <YOU@YOUR_DOMAIN_NAME>",
		"Hey %recipient.first%",
		"If you wish to unsubscribe, click http://mailgun/unsubscribe/%recipient.id%",
	) // IMPORTANT: No To:-field recipients!

	m.AddRecipientAndVariables("bob@example.com", map[string]interface{}{
		"first": "bob",
		"id":    1,
	})

	m.AddRecipientAndVariables("alice@example.com", map[string]interface{}{
		"first": "alice",
		"id":    2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func UpdateDomainConnection(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateDomainConnection(ctx, domain, mailgun.DomainConnection{
		RequireTLS:       true,
		SkipVerification: true,
	})
}

func UpdateMember(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := mg.UpdateMember(ctx, "bar@example.com", "list@example.com", mailgun.Member{
		Name:       "Foo Bar",
		Subscribed: mailgun.Unsubscribed,
	})
	return err
}

func UpdateWebhook(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateWebhook(ctx, "clicked", []string{"https://your_domain.com/clicked"})
}

func VerifyWebhookSignature(domain, apiKey, timestamp, token, signature string) (bool, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	return mg.VerifyWebhookSignature(mailgun.Signature{
		TimeStamp: timestamp,
		Token:     token,
		Signature: signature,
	})
}

func SendMessageWithTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create a new template
	err = mg.CreateTemplate(ctx, &mailgun.Template{
		Name: "my-template",
		Version: mailgun.TemplateVersion{
			Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
			Engine:   mailgun.TemplateEngineGo,
			Tag:      "v1",
		},
	})
	if err != nil {
		return err
	}

	// Give time for template to show up in the system.
	time.Sleep(time.Second * 1)

	// Create a new message with template
	m := mg.NewMessage("Excited User <excited@example.com>", "Template example", "")
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
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateTemplate(ctx, &mailgun.Template{
		Name: "my-template",
		Version: mailgun.TemplateVersion{
			Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
			Engine:   mailgun.TemplateEngineGo,
			Tag:      "v1",
		},
	})
}

func DeleteTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.DeleteTemplate(ctx, "my-template")
}

func UpdateTemplate(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateTemplate(ctx, &mailgun.Template{
		Name:        "my-template",
		Description: "Add a description to the template",
	})
}

func GetTemplate(domain, apiKey string) (mailgun.Template, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetTemplate(ctx, "my-template")
}

func ListActiveTemplates(domain, apiKey string) ([]mailgun.Template, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListTemplates(&mailgun.ListTemplateOptions{Active: true})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Template
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ListTemplates(domain, apiKey string) ([]mailgun.Template, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListTemplates(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Template
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func AddTemplateVersion(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.AddTemplateVersion(ctx, "my-template", &mailgun.TemplateVersion{
		Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
		Engine:   mailgun.TemplateEngineGo,
		Tag:      "v2",
		Active:   true,
	})
}

func DeleteTemplateVersion(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Delete the template version tagged as 'v2'
	return mg.DeleteTemplateVersion(ctx, "my-template", "v2")
}

func GetTemplateVersion(domain, apiKey string) (mailgun.TemplateVersion, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Get the template version tagged as 'v2'
	return mg.GetTemplateVersion(ctx, "my-template", "v2")
}

func UpdateTemplateVersion(domain, apiKey string) error {
	mg := mailgun.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.UpdateTemplateVersion(ctx, "my-template", &mailgun.TemplateVersion{
		Comment: "Add a comment to the template and make it 'active'",
		Tag:     "v2",
		Active:  true,
	})
}

func ListTemplateVersions(domain, apiKey string) ([]mailgun.TemplateVersion, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListTemplateVersions("my-template", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.TemplateVersion
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}
