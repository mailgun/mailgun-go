package mailgun

import (
	"log"

	"github.com/facebookgo/ensure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListEvents()", func() {
	var t GinkgoTInterface
	var it *EventIterator
	var mg Mailgun
	var err error

	BeforeEach(func() {
		mg, err = NewMailgunFromEnv()
		Expect(err).To(BeNil())
		it = mg.ListEvents(&EventsOptions{Limit: 5})
	})

	Describe("it.Next()", func() {
		It("Should iterate forward through pages of events", func() {
			var firstPage, secondPage []Event

			ensure.True(t, it.Next(&firstPage))
			ensure.True(t, it.NextURL != "")
			ensure.True(t, len(firstPage) != 0)
			firstIterator := *it

			ensure.True(t, it.Next(&secondPage))
			ensure.True(t, len(secondPage) != 0)

			// Pages should be different
			ensure.NotDeepEqual(t, firstPage, secondPage)
			ensure.True(t, firstIterator.NextURL != it.NextURL)
			ensure.True(t, firstIterator.PrevURL != it.PrevURL)
			ensure.Nil(t, it.Err())
		})
	})

	Describe("it.Previous()", func() {
		It("Should iterate backward through pages of events", func() {
			var firstPage, secondPage, previousPage []Event
			ensure.True(t, it.Next(&firstPage))
			ensure.True(t, it.Next(&secondPage))

			ensure.True(t, it.Previous(&previousPage))
			ensure.True(t, len(previousPage) != 0)
			ensure.DeepEqual(t, previousPage, firstPage)
		})
	})

	Describe("it.First()", func() {
		It("Should retrieve the first page of events", func() {
			var firstPage, secondPage []Event
			ensure.True(t, it.First(&firstPage))
			ensure.True(t, it.First(&secondPage))
			ensure.DeepEqual(t, firstPage, secondPage)

			// Calling first resets the iterator to the first page
			ensure.True(t, it.Next(&secondPage))
			ensure.NotDeepEqual(t, firstPage, secondPage)
		})
	})

	Describe("it.Last()", func() {
		It("Should retrieve the last page of events", func() {
			var firstPage, lastPage, previousPage []Event
			// Calling Last() is invalid unless you first use First() or Next()
			ensure.False(t, it.Last(&lastPage))
			ensure.True(t, len(lastPage) == 0)

			ensure.True(t, it.Next(&firstPage))
			ensure.True(t, len(firstPage) != 0)

			ensure.True(t, it.Last(&lastPage))
			ensure.True(t, len(lastPage) != 0)

			// Calling first resets the iterator to the first page
			ensure.True(t, it.Previous(&previousPage))
			ensure.NotDeepEqual(t, lastPage, previousPage)
		})
	})
})

var _ = Describe("EventIterator()", func() {
	log := log.New(GinkgoWriter, "EventIterator() - ", 0)
	var t GinkgoTInterface
	var mg Mailgun
	var err error

	BeforeEach(func() {
		t = GinkgoT()
		mg, err = NewMailgunFromEnv()
		ensure.Nil(t, err)
	})

	Describe("GetFirstPage()", func() {
		Context("When no parameters are supplied", func() {
			It("Should return a list of events", func() {
				ei := mg.NewEventIterator()
				err := ei.GetFirstPage(GetEventsOptions{})
				ensure.Nil(t, err)

				// Print out the kind of event and timestamp.
				// Specifics about each event will depend on the "event" type.
				events := ei.Events()
				log.Printf("Event\tTimestamp\t")
				for _, event := range events {
					log.Printf("%s\t%v\t\n", event["event"], event["timestamp"])
				}
				log.Printf("%d events dumped\n\n", len(events))
				ensure.True(t, len(events) != 0)

				// TODO: (thrawn01) The more I look at this and test it,
				// the more I doubt it will ever work consistently
				//ei.GetPrevious()
			})
		})
	})
})
