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
	"github.com/mailgun/mailgun-go/v4/events"
	"github.com/mailgun/mailgun-go/v4/mtypes"
)

// TODO(v5): remove/move?
type MockServer interface {
	Stop()
	URL() string
	DomainIPS() []string
	DomainList() []DomainContainer
	ExportList() []mtypes.Export
	MailingList() []MailingListContainer
	RouteList() []mtypes.Route
	Events() []events.Event
	Webhooks() mtypes.WebHooksListResponse
	Templates() []mtypes.Template
	SubaccountList() []mtypes.Subaccount
}

// A mailgun api mock suitable for testing
// TODO(v5): rename to MockServer?
type mockServer struct {
	srv *httptest.Server

	domainIPS        []string
	domainList       []DomainContainer
	exportList       []mtypes.Export
	mailingList      []MailingListContainer
	routeList        []mtypes.Route
	events           []events.Event
	templates        []mtypes.Template
	templateVersions map[string][]mtypes.TemplateVersion
	unsubscribes     []mtypes.Unsubscribe
	complaints       []mtypes.Complaint
	bounces          []mtypes.Bounce
	credentials      []mtypes.Credential
	tags             []mtypes.Tag
	subaccountList   []mtypes.Subaccount
	webhooks         mtypes.WebHooksListResponse
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

func (ms *mockServer) ExportList() []mtypes.Export {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.exportList
}

func (ms *mockServer) MailingList() []MailingListContainer {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.mailingList
}

func (ms *mockServer) RouteList() []mtypes.Route {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.routeList
}

func (ms *mockServer) Events() []events.Event {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.events
}

func (ms *mockServer) Webhooks() mtypes.WebHooksListResponse {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.webhooks
}

func (ms *mockServer) Templates() []mtypes.Template {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.templates
}

func (ms *mockServer) Unsubscribes() []mtypes.Unsubscribe {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.unsubscribes
}

func (ms *mockServer) SubaccountList() []mtypes.Subaccount {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.subaccountList
}

// Create a new instance of the mailgun API mock server
func NewMockServer() MockServer {
	ms := mockServer{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		ms.addIPRoutes(r)
		ms.addExportRoutes(r)
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
		ms.addTagsRoutes(r)
	})
	r.Route("/v5", func(r chi.Router) {
		ms.addSubaccountRoutes(r)
	})
	ms.addDomainRoutes(r) // mix of v3 and v4
	ms.addValidationRoutes(r)
	ms.addAnalyticsRoutes(r)

	// Start the server
	ms.srv = httptest.NewServer(r)
	return &ms
}

// Stop the server
func (ms *mockServer) Stop() {
	ms.srv.Close()
}

func (ms *mockServer) URL() string {
	return ms.srv.URL
}

func toJSON(w http.ResponseWriter, obj any) {
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

func stringToMap(v string) map[string]any {
	if v == "" {
		return nil
	}

	result := make(map[string]any)
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

type okResp struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}
