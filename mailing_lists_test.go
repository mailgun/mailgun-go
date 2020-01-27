package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
)

func TestMailingListMembers(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	address := randomEmail("list", testDomain)
	_, err := mg.CreateMailingList(ctx, mailgun.MailingList{
		Address:     address,
		Name:        address,
		Description: "TestMailingListMembers-related mailing list",
		AccessLevel: mailgun.AccessLevelMembers,
	})
	ensure.Nil(t, err)
	defer func() {
		ensure.Nil(t, mg.DeleteMailingList(ctx, address))
	}()

	var countMembers = func() int {
		var page []mailgun.Member
		var count int

		it := mg.ListMembers(address, nil)
		for it.Next(ctx, &page) {
			count += len(page)
		}
		ensure.Nil(t, it.Err())
		return count
	}

	startCount := countMembers()
	protoJoe := mailgun.Member{
		Address:    "joe@example.com",
		Name:       "Joe Example",
		Subscribed: mailgun.Subscribed,
	}
	ensure.Nil(t, mg.CreateMember(ctx, true, address, protoJoe))
	newCount := countMembers()
	ensure.False(t, newCount <= startCount)

	theMember, err := mg.GetMember(ctx, "joe@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Address, protoJoe.Address)
	ensure.DeepEqual(t, theMember.Name, protoJoe.Name)
	ensure.DeepEqual(t, theMember.Subscribed, protoJoe.Subscribed)
	ensure.True(t, len(theMember.Vars) == 0)

	_, err = mg.UpdateMember(ctx, "joe@example.com", address, mailgun.Member{
		Name: "Joe Cool",
	})
	ensure.Nil(t, err)

	theMember, err = mg.GetMember(ctx, "joe@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Name, "Joe Cool")
	ensure.Nil(t, mg.DeleteMember(ctx, "joe@example.com", address))
	ensure.DeepEqual(t, countMembers(), startCount)

	err = mg.CreateMemberList(ctx, nil, address, []interface{}{
		mailgun.Member{
			Address:    "joe.user1@example.com",
			Name:       "Joe's debugging account",
			Subscribed: mailgun.Unsubscribed,
		},
		mailgun.Member{
			Address:    "Joe Cool <joe.user2@example.com>",
			Name:       "Joe's Cool Account",
			Subscribed: mailgun.Subscribed,
		},
		mailgun.Member{
			Address: "joe.user3@example.com",
			Vars: map[string]interface{}{
				"packet-email": "KW9ABC @ BOGBBS-4.#NCA.CA.USA.NOAM",
			},
		},
	})
	ensure.Nil(t, err)

	theMember, err = mg.GetMember(ctx, "joe.user2@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Name, "Joe's Cool Account")
	ensure.NotNil(t, theMember.Subscribed)
	ensure.True(t, *theMember.Subscribed)
}

func TestMailingLists(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	address := randomEmail("list", testDomain)
	protoList := mailgun.MailingList{
		Address:     address,
		Name:        "List1",
		Description: "A list created by an acceptance test.",
		AccessLevel: mailgun.AccessLevelMembers,
	}

	var countLists = func() int {
		var count int
		it := mg.ListMailingLists(nil)
		var page []mailgun.MailingList
		for it.Next(ctx, &page) {
			count += len(page)
		}
		ensure.Nil(t, it.Err())
		return count
	}

	_, err := mg.CreateMailingList(ctx, protoList)
	ensure.Nil(t, err)
	defer func() {
		ensure.Nil(t, mg.DeleteMailingList(ctx, address))

		_, err := mg.GetMailingList(ctx, address)
		ensure.NotNil(t, err)
	}()

	actualCount := countLists()
	ensure.False(t, actualCount < 1)

	theList, err := mg.GetMailingList(ctx, address)
	ensure.Nil(t, err)

	protoList.CreatedAt = theList.CreatedAt // ignore this field when comparing.
	ensure.DeepEqual(t, theList, protoList)

	_, err = mg.UpdateMailingList(ctx, address, mailgun.MailingList{
		Description: "A list whose description changed",
	})
	ensure.Nil(t, err)

	theList, err = mg.GetMailingList(ctx, address)
	ensure.Nil(t, err)

	newList := protoList
	newList.Description = "A list whose description changed"
	ensure.DeepEqual(t, theList, newList)
}
