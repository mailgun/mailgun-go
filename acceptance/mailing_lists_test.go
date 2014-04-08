// +build acceptance,spendMoney

package acceptance

import (
	"testing"
	mailgun "github.com/mailgun/mailgun-go"
	"fmt"
)

func TestMailingLists(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	listAddr := fmt.Sprintf("list2@%s", domain)
	protoList := mailgun.List{
		Address: listAddr,
		Name: "List1",
		Description: "A list created by an acceptance test.",
		AccessLevel: mailgun.Members,
	}

	var countLists = func() int {
		total, _, err := mg.GetLists(mailgun.DefaultLimit, mailgun.DefaultSkip, "")
		if err != nil {
			t.Fatal(err)
		}
		return total
	}

	startCount := countLists()

	_, err := mg.CreateList(protoList)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = mg.DeleteList(listAddr)
		if err != nil {
			t.Fatal(err)
		}

		newCount := countLists()
		if newCount != startCount {
			t.Fatalf("Expected %d routes defined; got %d", startCount, newCount)
		}
	}()

	newCount := countLists()
	if newCount <= startCount {
		t.Fatalf("Expected %d routes defined; got %d", startCount+1, newCount)
	}

	theList, err := mg.GetListByAddress(listAddr)
	if err != nil {
		t.Fatal(err)
	}
	protoList.CreatedAt = theList.CreatedAt		// ignore this field when comparing.
	if theList != protoList {
		t.Fatalf("Unexpected list descriptor: Expected [%#v], Got [%#v]", protoList, theList)
	}

	_, err = mg.UpdateList(listAddr, mailgun.List{
		Description: "A list whose description changed",
	})
	if err != nil {
		t.Fatal(err)
	}

	theList, err = mg.GetListByAddress(listAddr)
	if err != nil {
		t.Fatal(err)
	}
	newList := protoList
	newList.Description = "A list whose description changed"
	if theList != newList {
		t.Fatalf("Expected [%#v], Got [%#v]", newList, theList)
	}
}