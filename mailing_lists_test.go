package mailgun_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

func TestPaging(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	it := mg.ListMailingLists(&mailgun.ListOptions{Limit: 10})
	var page []mailgun.MailingList
	var first []mailgun.MailingList

	/*for it.Next(&page) {
		spew.Dump(page)
	}*/
	//var count int
	it.Next(&first)
	it.First(&page)
	spew.Dump(page)
	//ensure.DeepEqual(t, page, first)
	/*if count > 2 {
		break
	}
	count++*/
	//}
	//ensure.Nil(t, it.Err())
}

/*func setup(t *testing.T) (Mailgun, string) {
	domain := reqEnv(t, "MG_DOMAIN")
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	address := fmt.Sprintf("%s@%s", strings.ToLower(randomString(6, "list")), domain)
	_, err = mg.CreateMailingList(MailingList{
		Address:     address,
		Name:        address,
		Description: "TestMailingListMembers-related mailing list",
		AccessLevel: Members,
	})
	ensure.Nil(t, err)
	return mg, address
}

func teardown(t *testing.T, mg Mailgun, address string) {
	ensure.Nil(t, mg.DeleteMailingList(address))
}

func TestMailingListMembers(t *testing.T) {
	mg, address := setup(t)
	defer teardown(t, mg, address)

	var countPeople = func() int {
		var count int
		it := mg.ListMembers(address, nil)
		var page []Member
		for it.Next(&page) {
			count += len(page)
		}
		ensure.Nil(t, it.Err())
		return count
	}

	startCount := countPeople()
	protoJoe := Member{
		Address:    "joe@example.com",
		Name:       "Joe Example",
		Subscribed: Subscribed,
	}
	ensure.Nil(t, mg.CreateMember(true, address, protoJoe))
	newCount := countPeople()
	ensure.False(t, newCount <= startCount)

	theMember, err := mg.GetMemberByAddress("joe@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Address, protoJoe.Address)
	ensure.DeepEqual(t, theMember.Name, protoJoe.Name)
	ensure.DeepEqual(t, theMember.Subscribed, protoJoe.Subscribed)
	ensure.True(t, len(theMember.Vars) == 0)

	_, err = mg.UpdateMember("joe@example.com", address, Member{
		Name: "Joe Cool",
	})
	ensure.Nil(t, err)

	theMember, err = mg.GetMemberByAddress("joe@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Name, "Joe Cool")
	ensure.Nil(t, mg.DeleteMember("joe@example.com", address))
	ensure.DeepEqual(t, countPeople(), startCount)

	err = mg.CreateMemberList(nil, address, []interface{}{
		Member{
			Address:    "joe.user1@example.com",
			Name:       "Joe's debugging account",
			Subscribed: Unsubscribed,
		},
		Member{
			Address:    "Joe Cool <joe.user2@example.com>",
			Name:       "Joe's Cool Account",
			Subscribed: Subscribed,
		},
		Member{
			Address: "joe.user3@example.com",
			Vars: map[string]interface{}{
				"packet-email": "KW9ABC @ BOGBBS-4.#NCA.CA.USA.NOAM",
			},
		},
	})
	ensure.Nil(t, err)

	theMember, err = mg.GetMemberByAddress("joe.user2@example.com", address)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, theMember.Name, "Joe's Cool Account")
	ensure.NotNil(t, theMember.Subscribed)
	ensure.True(t, *theMember.Subscribed)
}

func TestMailingLists(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	listAddr := fmt.Sprintf("%s@%s", strings.ToLower(randomString(7, "list")), domain)
	protoList := MailingList{
		Address:     listAddr,
		Name:        "List1",
		Description: "A list created by an acceptance test.",
		AccessLevel: Members,
	}

	var countLists = func() int {
		var count int
		it := mg.ListMailingLists(nil)
		var page []MailingList
		for it.Next(&page) {
			count += len(page)
		}
		ensure.Nil(t, it.Err())
		return count
	}

	_, err = mg.CreateMailingList(protoList)
	ensure.Nil(t, err)
	defer func() {
		ensure.Nil(t, mg.DeleteMailingList(listAddr))

		_, err := mg.GetMailingList(listAddr)
		ensure.NotNil(t, err)
	}()

	actualCount := countLists()
	ensure.False(t, actualCount < 1)

	theList, err := mg.GetMailingList(listAddr)
	ensure.Nil(t, err)

	protoList.CreatedAt = theList.CreatedAt // ignore this field when comparing.
	ensure.DeepEqual(t, theList, protoList)

	_, err = mg.UpdateMailingList(listAddr, MailingList{
		Description: "A list whose description changed",
	})
	ensure.Nil(t, err)

	theList, err = mg.GetMailingList(listAddr)
	ensure.Nil(t, err)

	newList := protoList
	newList.Description = "A list whose description changed"
	ensure.DeepEqual(t, theList, newList)
}*/
