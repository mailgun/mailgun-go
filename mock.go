package mailgun

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type MockServer interface {
	Stop()
	URL4() string
	URL() string
	DomainIPS() []string
	DomainList() []DomainContainer
	ExportList() []Export
	MailingList() []MailingListContainer
	RouteList() []Route
	Events() []Event
	Webhooks() WebHooksListResponse
	Templates() []Template
}

// A mailgun api mock suitable for testing
type mockServer struct {
	srv *httptest.Server

	domainIPS        []string
	domainList       []DomainContainer
	exportList       []Export
	mailingList      []MailingListContainer
	routeList        []Route
	events           []Event
	templates        []Template
	templateVersions map[string][]TemplateVersion
	unsubscribes     []Unsubscribe
	complaints       []Complaint
	bounces          []Bounce
	credentials      []Credential
	stats            []Stats
	tags             []Tag
	webhooks         WebHooksListResponse
	mutex            sync.Mutex
}

func (ms *mockServer) DomainIPS() []string {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.domainIPS
}

func (ms *mockServer) DomainList() []DomainContainer {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.domainList
}

func (ms *mockServer) ExportList() []Export {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.exportList
}

func (ms *mockServer) MailingList() []MailingListContainer {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.mailingList
}

func (ms *mockServer) RouteList() []Route {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.routeList
}

func (ms *mockServer) Events() []Event {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.events
}

func (ms *mockServer) Webhooks() WebHooksListResponse {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.webhooks
}

func (ms *mockServer) Templates() []Template {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.templates
}

func (ms *mockServer) Unsubscribes() []Unsubscribe {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.unsubscribes
}

// Create a new instance of the mailgun API mock server
func NewMockServer() MockServer {
	ms := mockServer{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		ms.addIPRoutes(r)
		ms.addExportRoutes(r)
		ms.addDomainRoutes(r)
		ms.addMailingListRoutes(r)
		ms.addEventRoutes(r)
		ms.addMessagesRoutes(r)
		ms.addRoutes(r)
		ms.addWebhookRoutes(r)
		ms.addTemplateRoutes(r)
		ms.addTemplateVersionRoutes(r)
		ms.addUnsubscribesRoutes(r)
		ms.addComplaintsRoutes(r)
		ms.addBouncesRoutes(r)
		ms.addCredentialsRoutes(r)
		ms.addStatsRoutes(r)
		ms.addTagsRoutes(r)
	})
	ms.addValidationRoutes(r)

	// Start the server
	ms.srv = httptest.NewServer(r)
	return &ms
}

// Stop the server
func (ms *mockServer) Stop() {
	ms.srv.Close()
}

func (ms *mockServer) URL4() string {
	return ms.srv.URL + "/v4"
}

// URL returns the URL used to connect to the mock server
func (ms *mockServer) URL() string {
	return ms.srv.URL + "/v3"
}

func toJSON(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

func stringToBool(v string) bool {
	lower := strings.ToLower(v)
	if lower == "yes" || lower == "no" {
		return lower == "yes"
	}

	if v == "" {
		return false
	}

	result, err := strconv.ParseBool(v)
	if err != nil {
		panic(err)
	}
	return result
}

func stringToInt(v string) int {
	if v == "" {
		return 0
	}

	result, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(result)
}

func stringToMap(v string) map[string]interface{} {
	if v == "" {
		return nil
	}

	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(v), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func parseAddress(v string) string {
	if v == "" {
		return ""
	}
	e, err := mail.ParseAddress(v)
	if err != nil {
		panic(err)
	}
	return e.Address
}

// Given the page direction, pivot value and limit, calculate the offsets for the slice
func pageOffsets(pivotIdx []string, pivotDir, pivotVal string, limit int) (int, int) {
	switch pivotDir {
	case "first":
		if limit < len(pivotIdx) {
			return 0, limit
		}
		return 0, len(pivotIdx)
	case "last":
		if limit < len(pivotIdx) {
			return len(pivotIdx) - limit, len(pivotIdx)
		}
		return 0, len(pivotIdx)
	case "next":
		for i, item := range pivotIdx {
			if item == pivotVal {
				offset := i + 1 + limit
				if offset > len(pivotIdx) {
					offset = len(pivotIdx)
				}
				return i + 1, offset
			}
		}
		return 0, 0
	case "prev":
		for i, item := range pivotIdx {
			if item == pivotVal {
				if i == 0 {
					return 0, 0
				}

				offset := i - limit
				if offset < 0 {
					offset = 0
				}
				return offset, i
			}
		}
		return 0, 0
	}

	if limit > len(pivotIdx) {
		return 0, len(pivotIdx)
	}
	return 0, limit
}

func getPageURL(r *http.Request, params url.Values) string {
	if r.FormValue("limit") != "" {
		params.Add("limit", r.FormValue("limit"))
	}
	return "http://" + r.Host + r.URL.EscapedPath() + "?" + params.Encode()
}

// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func randomString(n int, prefix string) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}

func randomEmail(prefix, domain string) string {
	return strings.ToLower(fmt.Sprintf("%s@%s", randomString(20, prefix), domain))
}
