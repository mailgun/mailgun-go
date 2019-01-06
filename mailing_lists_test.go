package mailgun_test

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

func TestMailingListMembers(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	address := randomEmail("list", testDomain)
	_, err = mg.CreateMailingList(mailgun.MailingList{
		Address:     address,
		Name:        address,
		Description: "TestMailingListMembers-related mailing list",
		AccessLevel: mailgun.Members,
	})
	ensure.Nil(t, err)
	defer func() {
		ensure.Nil(t, mg.DeleteMailingList(address))
	}()

	var page []mailgun.MailingList
	it := mg.ListMailingLists(nil)
	for it.Next(&page) {
		ensure.DeepEqual(t, len(page) != 0, true)
	}

	var countMembers = func() int {
		var page []mailgun.Member
		var count int

		it := mg.ListMembers(address, nil)
		for it.Next(&page) {
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
	ensure.Nil(t, mg.CreateMember(true, address, protoJoe))
	newCount := countMembers()
	ensure.False(t, newCount <= startCount)

	theMember, err := mg.GetMember("joe@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Address, protoJoe.Address)
	ensure.DeepEqual(t, theMember.Name, protoJoe.Name)
	ensure.DeepEqual(t, theMember.Subscribed, protoJoe.Subscribed)
	ensure.True(t, len(theMember.Vars) == 0)

	_, err = mg.UpdateMember("joe@example.com", address, mailgun.Member{
		Name: "Joe Cool",
	})
	ensure.Nil(t, err)

	theMember, err = mg.GetMember("joe@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Name, "Joe Cool")
	ensure.Nil(t, mg.DeleteMember("joe@example.com", address))
	ensure.DeepEqual(t, countMembers(), startCount)

	err = mg.CreateMemberList(nil, address, []interface{}{
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

	theMember, err = mg.GetMember("joe.user2@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Name, "Joe's Cool Account")
	ensure.NotNil(t, theMember.Subscribed)
	ensure.True(t, *theMember.Subscribed)
}

func TestMailingLists(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	address := randomEmail("list", testDomain)
	protoList := mailgun.MailingList{
		Address:     address,
		Name:        "List1",
		Description: "A list created by an acceptance test.",
		AccessLevel: mailgun.Members,
	}

	var countLists = func() int {
		var count int
		it := mg.ListMailingLists(nil)
		var page []mailgun.MailingList
		for it.Next(&page) {
			count += len(page)
		}
		ensure.Nil(t, it.Err())
		return count
	}

	_, err = mg.CreateMailingList(protoList)
	ensure.Nil(t, err)
	defer func() {
		ensure.Nil(t, mg.DeleteMailingList(address))

		_, err := mg.GetMailingList(address)
		ensure.NotNil(t, err)
	}()

	actualCount := countLists()
	ensure.False(t, actualCount < 1)

	theList, err := mg.GetMailingList(address)
	ensure.Nil(t, err)

	protoList.CreatedAt = theList.CreatedAt // ignore this field when comparing.
	ensure.DeepEqual(t, theList, protoList)

	_, err = mg.UpdateMailingList(address, mailgun.MailingList{
		Description: "A list whose description changed",
	})
	ensure.Nil(t, err)

	theList, err = mg.GetMailingList(address)
	ensure.Nil(t, err)

	newList := protoList
	newList.Description = "A list whose description changed"
	ensure.DeepEqual(t, theList, newList)
}
