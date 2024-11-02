package mailgun_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mg.DeleteMailingList(ctx, address))
	}()

	var countMembers = func() int {
		var page []mailgun.Member
		var count int

		it := mg.ListMembers(address, nil)
		for it.Next(ctx, &page) {
			count += len(page)
		}
		require.NoError(t, it.Err())
		return count
	}

	startCount := countMembers()
	protoJoe := mailgun.Member{
		Address:    "joe@example.com",
		Name:       "Joe Example",
		Subscribed: mailgun.Subscribed,
	}
	require.NoError(t, mg.CreateMember(ctx, true, address, protoJoe))
	newCount := countMembers()
	require.False(t, newCount <= startCount)

	theMember, err := mg.GetMember(ctx, "joe@example.com", address)
	require.NoError(t, err)
	assert.Equal(t, protoJoe.Address, theMember.Address)
	assert.Equal(t, protoJoe.Name, theMember.Name)
	assert.Equal(t, protoJoe.Subscribed, theMember.Subscribed)
	assert.Len(t, theMember.Vars, 0)

	_, err = mg.UpdateMember(ctx, "joe@example.com", address, mailgun.Member{
		Name: "Joe Cool",
	})
	require.NoError(t, err)

	theMember, err = mg.GetMember(ctx, "joe@example.com", address)
	require.NoError(t, err)
	assert.Equal(t, "Joe Cool", theMember.Name)
	require.NoError(t, mg.DeleteMember(ctx, "joe@example.com", address))
	assert.Equal(t, startCount, countMembers())

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
	require.NoError(t, err)

	theMember, err = mg.GetMember(ctx, "joe.user2@example.com", address)
	require.NoError(t, err)
	assert.Equal(t, "Joe's Cool Account", theMember.Name)
	require.NotNil(t, theMember.Subscribed)
	assert.True(t, *theMember.Subscribed)
}

func TestMailingLists(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	address := randomEmail("list", testDomain)
	protoList := mailgun.MailingList{
		Address:         address,
		Name:            "List1",
		Description:     "A list created by an acceptance test.",
		AccessLevel:     mailgun.AccessLevelMembers,
		ReplyPreference: mailgun.ReplyPreferenceSender,
	}

	var countLists = func() int {
		var count int
		it := mg.ListMailingLists(nil)
		var page []mailgun.MailingList
		for it.Next(ctx, &page) {
			count += len(page)
		}
		require.NoError(t, it.Err())
		return count
	}

	_, err := mg.CreateMailingList(ctx, protoList)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mg.DeleteMailingList(ctx, address))

		_, err := mg.GetMailingList(ctx, address)
		require.NotNil(t, err)
	}()

	actualCount := countLists()
	require.False(t, actualCount < 1)

	theList, err := mg.GetMailingList(ctx, address)
	require.NoError(t, err)

	protoList.CreatedAt = theList.CreatedAt // ignore this field when comparing.
	assert.Equal(t, theList, protoList)

	_, err = mg.UpdateMailingList(ctx, address, mailgun.MailingList{
		Description: "A list whose description changed",
	})
	require.NoError(t, err)

	theList, err = mg.GetMailingList(ctx, address)
	require.NoError(t, err)

	newList := protoList
	newList.Description = "A list whose description changed"
	assert.Equal(t, theList, newList)
}

func TestListMailingListRegression(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()
	address := "test@example.com"

	_, err := mg.CreateMailingList(ctx, mailgun.MailingList{
		Address:     address,
		Name:        "paging",
		Description: "Test paging",
	})
	require.NoError(t, err)

	for i := 0; i < 200; i++ {
		var vars map[string]interface{}
		if i == 5 {
			vars = map[string]interface{}{"has": "vars"}
		}

		err := mg.CreateMember(ctx, false, address, mailgun.Member{
			Address: fmt.Sprintf("%03d@example.com", i),
			Vars:    vars,
		})
		require.NoError(t, err)
	}

	it := mg.ListMembers(address, nil)

	var members []mailgun.Member
	var found int
	for it.Next(ctx, &members) {
		for _, m := range members {
			if m.Vars != nil {
				found++
			}
		}
	}
	require.NoError(t, it.Err())
	assert.Equal(t, 1, found)
}
