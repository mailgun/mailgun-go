package mailgun

import (
	"log"
	"time"

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
		t = GinkgoT()
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
			ensure.True(t, len(firstPage) != 0)

			// Calling first resets the iterator to the first page
			ensure.True(t, it.Next(&secondPage))
			ensure.NotDeepEqual(t, firstPage, secondPage)
		})
	})

	Describe("it.Last()", func() {
		Context("If First() or Next() was not called first", func() {
			It("Should fail with error", func() {
				var lastPage []Event
				// Calling Last() is invalid unless you first use First() or Next()
				ensure.False(t, it.Last(&lastPage))
				ensure.True(t, len(lastPage) == 0)
			})

		})
	})

	Describe("it.Last()", func() {
		It("Should retrieve the last page of events", func() {
			var firstPage, lastPage, previousPage []Event
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

var _ = Describe("Event{}", func() {
	var t GinkgoTInterface

	BeforeEach(func() {
		t = GinkgoT()
	})

	Describe("ParseTimeStamp()", func() {
		Context("When 'timestamp' exists and is valid", func() {
			It("Should parse the timestamp into time.Time{}", func() {
				event := Event{
					"timestamp": 1476380259.578017,
				}
				timestamp, err := event.ParseTimeStamp()
				ensure.Nil(t, err)
				ensure.DeepEqual(t, timestamp, time.Date(2016, 10, 13, 17, 37, 39,
					578017*int(time.Microsecond/time.Nanosecond), time.UTC))

				event = Event{
					"timestamp": 1377211256.096436,
				}
				timestamp, err = event.ParseTimeStamp()
				ensure.Nil(t, err)
				ensure.DeepEqual(t, timestamp, time.Date(2013, 8, 22, 22, 40, 56,
					96436*int(time.Microsecond/time.Nanosecond), time.UTC))
			})
		})
		Context("When 'timestamp' is missing", func() {
			It("Should return error", func() {
				event := Event{
					"blah": "",
				}
				_, err := event.ParseTimeStamp()
				ensure.NotNil(t, err)
				ensure.DeepEqual(t, err.Error(), "'timestamp' field not found in event")
			})
		})
		Context("When 'timestamp' is not a float64", func() {
			It("Should return error", func() {
				event := Event{
					"timestamp": "1476380259.578017",
				}
				_, err := event.ParseTimeStamp()
				ensure.NotNil(t, err)
				ensure.DeepEqual(t, err.Error(), "'timestamp' field not a float64")
			})
		})
	})
})

var _ = Describe("PollEvents()", func() {
	log := log.New(GinkgoWriter, "PollEvents() - ", 0)
	var t GinkgoTInterface
	var it *EventPoller
	var mg Mailgun
	var err error

	BeforeEach(func() {
		t = GinkgoT()
	})

	Describe("it.Poll()", func() {
		It("Should return events once the threshold age has expired", func() {
			mg, err = NewMailgunFromEnv()
			Expect(err).To(BeNil())

			// Very short poll interval
			it = mg.PollEvents(&EventsOptions{
				// Poll() returns after this threshold is met
				// or events older than this threshold appear
				ThresholdAge: time.Second * 10,
				// Only events with a timestamp after this date/time will be returned
				Begin: time.Now().Add(time.Second * -3),
				// How often we poll the api for new events
				PollInterval: time.Second * 4})

			// Send an email
			toUser := reqEnv(t, "MG_EMAIL_TO")
			m := mg.NewMessage(fromUser, exampleSubject, exampleText, toUser)
			msg, id, err := mg.Send(m)
			ensure.Nil(t, err)

			log.Printf("New Email: %s Id: %s\n", msg, id)

			// Wait for our email event to arrive
			var events []Event
			it.Poll(&events)

			var found bool
			// Log the events we received
			for _, event := range events {
				eventMsg, _ := event.ParseMessageId()
				timeStamp, _ := event.ParseTimeStamp()
				log.Printf("Event: %s <%s> - %s", eventMsg, event["event"], timeStamp)

				// If we find our accepted email event
				if id == ("<"+eventMsg+">") && event["event"] == "accepted" {
					found = true
				}
			}
			// Ensure we found our email
			ensure.Nil(t, it.Err())
			ensure.True(t, found)
		})
	})
})
